package main

import (
	"fmt"
	"strings"
	"strconv"
	"time"
	game "cardooo/game"
	server "cardooo/net"
)

var mainGame game.Game

func main() {	

	mainGame = game.InitGame(server.SendMsg, server.BroadcastMessage)
	go server.StartTCP(AddNewAgent, RemoveAgent, ClientCommand)

	fmt.Println("[Cardooo] Server Start...")
	for {
		time.Sleep(1000*1000*1000)
	}
}

func AddNewAgent(id int) {
	a := mainGame.AddNewAgent(id)
	info := fmt.Sprintf("%v,%v,%v,%v",a.Id, a.MapId, a.Pos.X, a.Pos.Y)
	server.SendMsg(id, 1, 10, info)
	ServerCommand(id, 1, 1001, "")
}

func RemoveAgent(id int) {
	ServerCommand(id, 1, 1002, "")
	mainGame.RemoveAgent(id)
}

func ServerCommand(id int, sys int, api int, msg string) {
	fmt.Printf("[Cardooo][ServerCommand] %v,%v,%v,%s\n", id, sys, api, msg)

	switch api {
		case 1001:// [S>C] 玩家登場
			v := mainGame.AgentMap[id]		
			sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)	
			server.BroadcastMessage(id, sys, api, sendMsg)
		case 1002:// [S>C] 玩家離開
			v := mainGame.AgentMap[id]		
			sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)	
			server.BroadcastMessage(id, sys, api, sendMsg)						
		default:	
	}
}

func ClientCommand(id int, sys int, api int, msg string) {
	fmt.Printf("[Cardooo][ClientCommand] %v,%v,%v,%s\n", id, sys, api, msg)

	switch api {
	case 10://取得玩家狀態
		v := mainGame.AgentMap[id]		
		sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)		
		server.SendMsg(id, sys, api, sendMsg)

	case 11://取得服務器GAME狀態
		sendMsg := ""
		sendMsg += fmt.Sprintf("1000,%v|",mainGame.MapId)
		for k, v := range mainGame.AgentMap {
			sendMsg += fmt.Sprintf("%v,%v,%v,%v|", k, v.MapId, v.Pos.X, v.Pos.Y)
		}
		server.SendMsg(id, sys, api, sendMsg)
	case 12://玩家order
		server.BroadcastMessage(id, sys, api, msg)
		orderStr := strings.Split(msg, ",") // >> id,order
		order, _ := strconv.Atoi(orderStr[1]) // string >> int
		mainGame.OnOrder(id, order)
	case 100://戰報
		sendMsg := ""
		sendMsg += fmt.Sprintf("WORKING")
		fmt.Println(sendMsg)
		server.SendMsg(id, sys, api, sendMsg)
	default:
		
	}
}
