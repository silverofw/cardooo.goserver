package main

import (
	"fmt"
	"strings"
	"strconv"
	"time"
	"cardooo/enum"
	game "cardooo/game"
	server "cardooo/net"
	battle "cardooo/battle"
)

var mainGame game.Game
var mainBattle battle.Battle

func main() {	

	mainGame = game.InitLobyGame(server.SendMsg, server.BroadcastMessage)	

	go server.StartTCP(AddNewAgent, RemoveAgent, ClientCommand)
	mainBattle = battle.InitBattle()	

	fmt.Println("[Cardooo] Server Start...")
	for {
		time.Sleep(1000*1000*1000)
	}
}

func AddNewAgent(id int) {
	a := mainGame.AddNewAgentById(id)
	info := fmt.Sprintf("%v,%v,%v,%v",a.Id, a.MapId, a.Pos.X, a.Pos.Y)
	server.SendMsg(id, 1, MainEvent.CSC_PLAYER_STATE, info)
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
			if v != nil {
				sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)	
				server.BroadcastMessage(id, sys, api, sendMsg)
			}else{
				fmt.Printf("[ERROR][MAIN] agent is missing! id:%v \n", id)
			}
		default:	
	}
}

func ClientCommand(id int, sys int, api int, msg string) {
	fmt.Printf("[Cardooo][ClientCommand] %v,%v,%v,%s\n", id, sys, api, msg)	

	switch api {
	case MainEvent.CSC_PLAYER_STATE://取得玩家狀態
		v := mainGame.AgentMap[id]		
		sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)		
		server.SendMsg(id, sys, api, sendMsg)

	case MainEvent.CSC_GAME_STATE://取得服務器GAME狀態
		sendMsg := ""
		sendMsg += fmt.Sprintf("1000,%v|",mainGame.MapId)
		for k, v := range mainGame.AgentMap {
			sendMsg += fmt.Sprintf("%v,%v,%v,%v|", k, v.MapId, v.Pos.X, v.Pos.Y)
		}
		server.SendMsg(id, sys, api, sendMsg)
	case MainEvent.CSC_PLAYER_ORDER://玩家order
		server.BroadcastMessage(id, sys, api, msg)
		orderStr := strings.Split(msg, ",") // >> id,order
		order, _ := strconv.Atoi(orderStr[1]) // string >> int
		mainGame.OnOrder(id, order)
	case MainEvent.CSC_BATTLE_REPORT://戰報
		sendMsg := mainBattle.Report(id)
		fmt.Println("[BATTLE REPORT]: " + sendMsg)
		server.SendMsg(id, sys, api, sendMsg)
	case MainEvent.CSC_BATTLE_UPDATE_TEAM://更新隊伍
		agentsStr := strings.Split(msg, "|")
		agents := []game.Agent{}
		for _,v := range agentsStr {
			fmt.Printf("[BATTLE] UPATE PLAYER TEAM :%v \n", v)
			if v== "" {
				continue
			}
			agentStr := strings.Split(v, ",")
			pixel, _ := strconv.Atoi(agentStr[0])
			x, _ := strconv.Atoi(agentStr[1])
			y, _ := strconv.Atoi(agentStr[2])

			xOffset := 2
			agent := game.Agent{
				Pixel: pixel, MapId: 0, Hp: 10, Face: 0, 
				Pos: game.Pos{ X: x + xOffset, Y: y, },
			}
			agents = append(agents, agent)
		}
		mainBattle.UpdatePlayerTeam(agents)
		fmt.Printf("[BATTLE] UPATE PLAYER TEAM FINISH\n")
	default:
		
	}
}
