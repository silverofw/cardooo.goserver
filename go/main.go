package main

import (
	"fmt"
	//"bufio"
	"cardooo/core"
	server "cardooo/core/net"
	cardooo "cardooo/log"
	"net"
)

//module/package name

func main() {
	cardooo.Print()
	core.Print()
	server.Print()

	// 創建 TCP 監聽器，監聽所有網卡上的 1024 端口
	listener, _ := net.Listen("tcp", ":1024")
	println("啟動伺服器...")

	for {
		// 持續監聽客戶端連線
		conn, err := listener.Accept()
		if err != nil{
			println(err)
			continue
		}

		// 當客戶端連接時，創建一個新的 go 協程處理該客戶端
		go ClientLogic(conn)
	}
}

func ClientLogic(conn net.Conn) {
	fmt.Println("Client connected: " + conn.RemoteAddr().String())
	defer conn.Close()

    // 循環接收客戶端發送的消息
    for {
        // 創建一個 1024 字節的緩衝區
        buf := make([]byte, 1024)

        // 從連接中讀取數據，直到客戶端斷開連接為止
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println(err)
            return
        }

        // 將客戶端發送的消息回傳給客戶端
        conn.Write(buf[:n])
    }
}
