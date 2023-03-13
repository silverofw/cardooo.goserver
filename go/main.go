package main

import (
	"fmt"
	"bytes"
	//"strconv"
	"cardooo/core"
	server "cardooo/core/net"
	cardooo "cardooo/log"
	"net"
)

type client struct {
	uid string
	ip string
	conn net.Conn
}

var clients []client

var msgSize int = 1024
var port string = ":1024"

func main() {
	cardooo.Print()
	core.Print()
	server.Print()

	// 創建 TCP 監聽器，監聽所有網卡上的 1024 端口
	listener, _ := net.Listen("tcp", port)
	println("[cardooo][v.0.2] Server Start...")

	for {
		// 持續監聽客戶端連線
		conn, err := listener.Accept()
		if err != nil{
			println(err)
			continue
		}

		newClient := client{
			conn: conn, 
			ip: conn.RemoteAddr().String(),
			uid: "NewClient",
		}
		clients = append(clients, newClient)
		fmt.Println("[newClient]" + newClient.ip)
		msg := string("Wellcome!NewClient!(" + newClient.ip + ")")
		newClient.sendToC("0001", "0001", msg)

		go newClient.handleClient()
	}
}


func (c *client)handleClient() {
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
	for i, client := range clients {
		if client.conn == c.conn {
			fmt.Println("[removeClient]" + client.conn.RemoteAddr().String())
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func (c *client)processMsg(buf []byte) {
	systemId := string(bytes.Trim(buf[4:8], "\x00"))
	apiId := string(bytes.Trim(buf[8:12], "\x00"))
	params := string(bytes.Trim(buf[12:], "\x00"))

	switch apiId {
	case "0001":
		//server well come new client~
	case "0002":
		c.setUid(systemId, apiId, params)
	case "0003":
		msg := fmt.Sprintf("[B-%v-%v]:%v", c.uid, c.ip, params) 
		broadcastMessage(systemId, apiId, msg)
	case "9999":
		c.sendToC(systemId, apiId, params)
	default:
		// 將客戶端發送的消息回傳給客戶端
		c.conn.Write(buf)
	}
}

func (c *client)setUid(systemId string, apiId string, params string) {
	c.uid = params
	fmt.Printf("[setUid]: %s\n", c.uid)
	msg := fmt.Sprintf("Hello! %s! Set uid finish!", c.uid)
	c.sendToC(systemId, apiId, msg)
}

func broadcastMessage(systemId string, apiId string, msg string) {	
	for _, c := range clients {	
		c.sendToC(systemId, apiId, msg)
	}
}

func (c *client)sendToC(systemId string, apiId string, params string) {	
	msg := fmt.Sprintf("%s%s%s", systemId, apiId, params)
	buf := []byte(msg)
	//fmt.Println(msg)
	_, err := c.conn.Write(buf)
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}
