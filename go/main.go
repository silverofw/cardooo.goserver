package main

import (
	"fmt"
	"time"
	game "cardooo/game"
	server "cardooo/net"
)

var mainGame game.Game

func main() {	

	mainGame = game.InitGame()
	mainGame.SendMsg = server.SendMsg	
	go server.StartTCP(AddNewAgent, RemoveAgent, ClientCommand)

	fmt.Println("[Cardooo] Server Start...")
	for {
		time.Sleep(1000*1000*1000)
	}
}

func AddNewAgent(id int) {
	mainGame.AddNewAgent(id)
}

func RemoveAgent(id int) {
	mainGame.RemoveAgent(id)
}

func ClientCommand(uid string, systemId string, apiId string, msg string) {
	fmt.Printf("[Cardooo][ClientCommand] %s,%s,%s,%s\n", uid, systemId, apiId, msg)

	switch apiId {
	case "0010"://取得服務器GAME狀態
	case "0011"://角色移動
	case "0012"://角色行為
	default:
		
	}
}
