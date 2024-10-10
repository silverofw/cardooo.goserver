package common

import (
	"fmt"
	"net"
)

type Client struct {
	Id      int
	Account int
	Name    string
	Ip      string
	Frame   int
	Conn    net.Conn
	IsConn  bool
}

func (c *Client) SendToClient(sys int, api int, msg string) {
	systemId := fmt.Sprintf("%04d", sys)
	apiId := fmt.Sprintf("%04d", api)
	c.SendToC(systemId, apiId, msg)
}

func (c *Client) SendToC(systemId string, apiId string, params string) {
	msg := fmt.Sprintf("[<]%s%s%s[>]", systemId, apiId, params)
	buf := []byte(msg)
	//fmt.Println(msg)
	_, err := c.Conn.Write(buf)
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}
func (c *Client) SetUid(systemId string, apiId string, params string) {
	c.Name = params
	fmt.Printf("[setUid]: %s\n", c.Name)
	msg := fmt.Sprintf("Hello! %s! Set uid finish!", c.Name)
	c.SendToC(systemId, apiId, msg)
}
