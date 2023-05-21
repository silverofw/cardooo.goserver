package server

import (
	"fmt"
	"net"
	"bytes"
)

type client struct {
	id int
	name string
	ip string
	frame int
	conn net.Conn
	isConn bool
}

var token int = 0
var clients map[int]client
var msgSize int = 1024
var port string = ":1024"

var newClient func(int)
var delClient func(int)
var clientCommand func(string, string, string, string) 

func StartTCP(newC func(int), delC func(int), clientC func(string, string, string, string)){
	newClient = newC
	delClient = delC
	clientCommand = clientC
	clients = make(map[int]client)

	// 創建 TCP 監聽器，監聽所有網卡上的 1024 端口
	listener, _ := net.Listen("tcp", port)

	for {
		// 持續監聽客戶端連線
		conn, err := listener.Accept()
		if err != nil{
			println(err)
			continue
		}
		
		token++
		newClient := client{
			id: token,
			conn: conn, 
			ip: conn.RemoteAddr().String(),
			frame: 0,
			name: "NewClient",
			isConn: true,
		}
		clients[newClient.id] = newClient
		msg := fmt.Sprintf("[Server][%v]Wellcome!NewClient!(%s)", newClient.id, newClient.ip)
		fmt.Println(msg)
		//msg := string("Wellcome!NewClient!(" + newClient.ip + ")")		
		newClient.sendToC("0001", "0001", msg)

		go newClient.handleClient()		
	}
}


func (c *client)handleClient() {
	newClient(c.id)

	for {
		buf := make([]byte, msgSize)
		_, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		//fmt.Println("Received message:", msg)
		c.processMsg(buf)
	}

	c.conn.Close()
	c.removeClient()
}

func (c *client)removeClient() {
	c.isConn = false
	delClient(c.id)
	fmt.Printf("[Server][removeClient][%v][ip: %s]\n",c.id, c.ip)	
	delete(clients, c.id)
}

func (c *client)processMsg(buf []byte) {
	uid := string(bytes.Trim(buf[0:4], "\x00"))
	systemId := string(bytes.Trim(buf[4:8], "\x00"))
	apiId := string(bytes.Trim(buf[8:12], "\x00"))
	params := string(bytes.Trim(buf[12:], "\x00"))

	switch apiId {
	case "0002":
		c.setUid(systemId, apiId, params)
	case "0003":
		msg := fmt.Sprintf("[%v] %s: %v", c.id, c.name, params) 
		BroadcastMessage(systemId, apiId, msg)
	case "9999":
		c.sendToC(systemId, apiId, params)
	default:
		// 將客戶端發送的消息回傳給客戶端
		c.conn.Write(buf)
		clientCommand(uid, systemId, apiId, params)
	}
}

func (c *client)setUid(systemId string, apiId string, params string) {
	c.name = params
	fmt.Printf("[setUid]: %s\n", c.name)
	msg := fmt.Sprintf("Hello! %s! Set uid finish!", c.name)
	c.sendToC(systemId, apiId, msg)
}

func BroadcastMessage(systemId string, apiId string, msg string) {	
	for _, c := range clients {	
		c.sendToC(systemId, apiId, msg)
	}
}

func SendMsg(id int, systemId int, apiId int, params string) {	
	sys := fmt.Sprintf("%04d",systemId)
	api := fmt.Sprintf("%04d",apiId)
	if id != 0 {
		c := clients[id]
		c.sendToC(sys, api, params)
	} else {
		BroadcastMessage(sys, api, params)
	}
}

func (c *client)sendToC(systemId string, apiId string, params string) {	
	msg := fmt.Sprintf("[<]%s%s%s[>]", systemId, apiId, params)
	buf := []byte(msg)
	//fmt.Println(msg)
	_, err := c.conn.Write(buf)
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}
