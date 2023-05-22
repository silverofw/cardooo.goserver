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

func ClientCommand(id int, sys int, api int, msg string) {
	fmt.Printf("[Cardooo][ClientCommand] %v,%v,%v,%s\n", id, sys, api, msg)

	switch api {
	case 10://取得服務器GAME狀態
		sendMsg := ""
		for k, v := range mainGame.AgentMap {
			sendMsg += fmt.Sprintf("%v,%v,%v|", k, v.Pos.X, v.Pos.Y)
		}
		server.SendMsg(id, sys, api, sendMsg)
	case 11://角色移動
	case 12://角色行為
	default:
		
	}
}
