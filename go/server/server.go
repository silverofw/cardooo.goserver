package server

import (
	"bytes"
	"cardooo/common"
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	token   int
	Clients map[int]*common.Client
	msgSize int
	port    string

	newClient     func(int)
	delClient     func(*common.Client)
	clientCommand func(int, int, int, string)
}

func InitServer(newC func(int), delC func(*common.Client), clientC func(int, int, int, string)) *Server {
	return &Server{
		token:   1000,
		msgSize: 1024,
		port:    ":1024",

		newClient:     newC,
		delClient:     delC,
		clientCommand: clientC,
		Clients:       make(map[int]*common.Client),
	}
}

func (s *Server) Stop() {
	fmt.Println("Stop server")
}

func (s *Server) StartTCP() {

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
		newClient := common.Client{
			Id:     s.token,
			Conn:   conn,
			Ip:     conn.RemoteAddr().String(),
			Frame:  0,
			Name:   "NewClient",
			IsConn: true,
		}
		s.Clients[newClient.Id] = &newClient
		msg := fmt.Sprintf("[Server][%v]Wellcome!NewClient!(%s)", newClient.Id, newClient.Ip)
		fmt.Println(msg)
		//msg := string("Wellcome!NewClient!(" + newClient.ip + ")")
		newClient.SendToC("0001", "0001", msg)

		go handleClient(s, &newClient)
	}
}

func (s *Server) UpdateAccount(id int, account int) {
	c := s.Clients[id]
	c.Account = account
	s.Clients[id] = c
}
func handleClient(s *Server, c *common.Client) {
	s.newClient(c.Id)

	for {
		buf := make([]byte, s.msgSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		//fmt.Println("Received message:", msg)
		processMsg(s, c, buf)
	}

	c.Conn.Close()
	removeClient(s, c)
}

func removeClient(s *Server, c *common.Client) {
	c.IsConn = false
	s.delClient(c)
	fmt.Printf("[Server][removeClient][%v][ip: %s]\n", c.Id, c.Ip)
	delete(s.Clients, c.Id)
}

func processMsg(s *Server, c *common.Client, buf []byte) {
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
		c.SetUid(systemStr, apiStr, params)
	case 3:
		msg := fmt.Sprintf("[%v] %s: %v", c.Id, c.Name, params)
		fmt.Printf("[api:3] %v\n", msg)
		s.BroadcastMessage(-1, sys, api, msg)
	case 9999:
		c.SendToC(systemStr, apiStr, params)
	default:
		// 將客戶端發送的消息回傳給客戶端
		c.Conn.Write(buf)
		s.clientCommand(id, sys, api, params)
	}
}

func (s *Server) BroadcastMessage(id int, sys int, api int, msg string) {
	systemId := fmt.Sprintf("%04d", sys)
	apiId := fmt.Sprintf("%04d", api)
	for _, c := range s.Clients {
		if c.Id != id {
			c.SendToC(systemId, apiId, msg)
		}
	}
}

func (s *Server) SendMsg(id int, sys int, api int, msg string) {
	c := s.Clients[id]
	systemId := fmt.Sprintf("%04d", sys)
	apiId := fmt.Sprintf("%04d", api)
	c.SendToC(systemId, apiId, msg)
}
