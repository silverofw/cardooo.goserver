package server

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type client struct {
	id      int
	Account int
	name    string
	ip      string
	frame   int
	conn    net.Conn
	isConn  bool
}

type Server struct {
	token   int
	Clients map[int]client
	msgSize int
	port    string

	newClient     func(int)
	delClient     func(int)
	clientCommand func(int, int, int, string)
}

func InitServer(newC func(int), delC func(int), clientC func(int, int, int, string)) Server {
	s := Server{
		token:   1000,
		msgSize: 1024,
		port:    ":1024",

		newClient:     newC,
		delClient:     delC,
		clientCommand: clientC,
		Clients:       make(map[int]client),
	}
	return s
}

func (s *Server) StartTCP(newC func(int), delC func(int), clientC func(int, int, int, string)) {

	// 創建 TCP 監聽器，監聽所有網卡上的 1024 端口
	listener, _ := net.Listen("tcp", s.port)

	for {
		// 持續監聽客戶端連線
		conn, err := listener.Accept()
		if err != nil {
			println(err)
			continue
		}

		s.token++
		newClient := client{
			id:     s.token,
			conn:   conn,
			ip:     conn.RemoteAddr().String(),
			frame:  0,
			name:   "NewClient",
			isConn: true,
		}
		s.Clients[newClient.id] = newClient
		msg := fmt.Sprintf("[Server][%v]Wellcome!NewClient!(%s)", newClient.id, newClient.ip)
		fmt.Println(msg)
		//msg := string("Wellcome!NewClient!(" + newClient.ip + ")")
		newClient.sendToC("0001", "0001", msg)

		go newClient.handleClient(s)
	}
}

func (c *client) handleClient(s *Server) {
	s.newClient(c.id)

	for {
		buf := make([]byte, s.msgSize)
		_, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		//fmt.Println("Received message:", msg)
		c.processMsg(s, buf)
	}

	c.conn.Close()
	c.removeClient(s)
}

func (s *Server) UpdateAccount(id int, account int) {
	c := s.Clients[id]
	c.Account = account
	s.Clients[id] = c
}

func (c *client) removeClient(s *Server) {
	c.isConn = false
	s.delClient(c.id)
	fmt.Printf("[Server][removeClient][%v][ip: %s]\n", c.id, c.ip)
	delete(s.Clients, c.id)
}

func (c *client) processMsg(s *Server, buf []byte) {
	idStr := string(bytes.Trim(buf[0:4], "\x00"))
	systemStr := string(bytes.Trim(buf[4:8], "\x00"))
	apiStr := string(bytes.Trim(buf[8:12], "\x00"))
	params := string(bytes.Trim(buf[12:], "\x00"))

	id, _ := strconv.Atoi(idStr)      // string >> int
	sys, _ := strconv.Atoi(systemStr) // string >> int
	api, _ := strconv.Atoi(apiStr)    // string >> int

	_, ok := s.Clients[id]
	if !ok {
		fmt.Printf("[ERROR][SERVER][%v,%v,%v] wrong id!\n", id, sys, api)
		return
	}

	switch api {
	case 2:
		c.setUid(systemStr, apiStr, params)
	case 3:
		msg := fmt.Sprintf("[%v] %s: %v", c.id, c.name, params)
		fmt.Printf("[api:3] %v\n", msg)
		s.BroadcastMessage(-1, sys, api, msg)
	case 9999:
		c.sendToC(systemStr, apiStr, params)
	default:
		// 將客戶端發送的消息回傳給客戶端
		c.conn.Write(buf)
		s.clientCommand(id, sys, api, params)
	}
}

func (c *client) setUid(systemId string, apiId string, params string) {
	c.name = params
	fmt.Printf("[setUid]: %s\n", c.name)
	msg := fmt.Sprintf("Hello! %s! Set uid finish!", c.name)
	c.sendToC(systemId, apiId, msg)
}

func (s *Server) BroadcastMessage(id int, sys int, api int, msg string) {
	systemId := fmt.Sprintf("%04d", sys)
	apiId := fmt.Sprintf("%04d", api)
	for _, c := range s.Clients {
		if c.id != id {
			c.sendToC(systemId, apiId, msg)
		}
	}
}

func (s *Server) SendMsg(id int, sys int, api int, msg string) {
	c := s.Clients[id]
	systemId := fmt.Sprintf("%04d", sys)
	apiId := fmt.Sprintf("%04d", api)
	c.sendToC(systemId, apiId, msg)
}

func (c *client) sendToC(systemId string, apiId string, params string) {
	msg := fmt.Sprintf("[<]%s%s%s[>]", systemId, apiId, params)
	buf := []byte(msg)
	//fmt.Println(msg)
	_, err := c.conn.Write(buf)
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}
