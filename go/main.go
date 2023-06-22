package main

import (
	"cardooo/battle"
	MainEvent "cardooo/enum"
	"cardooo/game"
	"cardooo/model"
	server "cardooo/net"
	"cardooo/user"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var userMgr user.UserMgr
var mainServer server.Server
var mainGame game.Game
var mainBattle battle.Battle

func main() {
	userMgr = user.InitUserMgr()
	mainServer = server.InitServer(AddNewAgent, RemoveAgent, ClientCommand)
	mainGame = game.InitLobyGame(mainServer.SendMsg, mainServer.BroadcastMessage)

	mainServer.StartTCP(AddNewAgent, RemoveAgent, ClientCommand)
	mainBattle = battle.InitBattle()

	fmt.Println("[Cardooo] Server Start...")
	for {
		time.Sleep(1000 * 1000 * 1000)
	}
}

func AddNewAgent(id int) {
	mainGame.AddNewAgentById(id)
	ClientCommand(id, 1, MainEvent.CSC_SERVER_STATE, "")
}

func RemoveAgent(id int) {
	ServerCommand(id, 1, 1002, "")
	mainGame.RemoveAgent(id)
}

func ServerCommand(id int, sys int, api int, msg string) {
	fmt.Printf("[Cardooo][ServerCommand] %v,%v,%v,%s\n", id, sys, api, msg)

	switch api {
	case 1001: // [S>C] 玩家登場
		v := mainGame.AgentMap[id]
		sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)
		mainServer.BroadcastMessage(id, sys, api, sendMsg)
	case 1002: // [S>C] 玩家離開
		v := mainGame.AgentMap[id]
		if v != nil {
			sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)
			mainServer.BroadcastMessage(id, sys, api, sendMsg)
		} else {
			fmt.Printf("[ERROR][MAIN] agent is missing! id:%v \n", id)
		}
	default:
	}
}

func ClientCommand(id int, sys int, api int, msg string) {
	fmt.Printf("[Cardooo][ClientCommand] %v,%v,%v,%s\n", id, sys, api, msg)

	switch api {
	case MainEvent.CSC_SERVER_STATE: //取得server狀態
		sendMsg := fmt.Sprintf("%v,%v", 1, id)
		mainServer.SendMsg(id, sys, api, sendMsg)
	case MainEvent.CSC_PLAYER_STATE: //取得玩家狀態
		orderStr := strings.Split(msg, ",")     // >> id,order
		account, _ := strconv.Atoi(orderStr[0]) // string >> int
		passward := orderStr[1]
		userMgr.UserLogin(account, passward, id)
		mainServer.UpdateAccount(id, account)

		v := mainGame.AgentMap[id]
		sendMsg := fmt.Sprintf("%v,%v,%v|", id, v.Pos.X, v.Pos.Y)
		u := userMgr.Users[mainServer.Clients[id].Account]
		for _, v := range u.Items {
			sendMsg += fmt.Sprintf("%v,%v,%v=", v.UID, v.Id, v.Quantity)
		}

		sendMsg += "|"
		for _, v := range u.Team {
			sendMsg += fmt.Sprintf("%v,%v,%v=", v.Pixel, v.Pos.X, v.Pos.Y)
		}

		mainServer.SendMsg(id, sys, api, sendMsg)
		ServerCommand(id, 1, 1001, "")

	case MainEvent.CSC_GAME_STATE: //取得服務器GAME狀態
		sendMsg := ""
		sendMsg += fmt.Sprintf("1000,%v|", mainGame.MapId)
		for k, v := range mainGame.AgentMap {
			sendMsg += fmt.Sprintf("%v,%v,%v,%v|", k, v.MapId, v.Pos.X, v.Pos.Y)
		}
		mainServer.SendMsg(id, sys, api, sendMsg)
	case MainEvent.CSC_PLAYER_ORDER: //玩家order
		mainServer.BroadcastMessage(id, sys, api, msg)
		orderStr := strings.Split(msg, ",")   // >> id,order
		order, _ := strconv.Atoi(orderStr[1]) // string >> int
		mainGame.OnOrder(id, order)
	case MainEvent.CSC_BATTLE_REPORT: //戰報
		orderStr := strings.Split(msg, ",")        // >> id,order
		battleType, _ := strconv.Atoi(orderStr[0]) // string >> int

		u := userMgr.Users[mainServer.Clients[id].Account]
		initData := battle.BattleInitData{
			PlayerTeam: u.Team,
		}

		switch battleType {
		case 0:
			fmt.Printf("[BATTLE REPORT] start normal battle!\n")
			initData.EnemyTeam = []model.Agent{
				{MapId: 0, Hp: 2, Face: 0, Pixel: 2001, Pos: model.Vector2{X: 1, Y: 2}},
				{MapId: 0, Hp: 2, Face: 0, Pixel: 2001, Pos: model.Vector2{X: 4, Y: 1}},
			}
			initData.EnemyName = "NORMAL_0001"
		case 1:
			enemyU := userMgr.GetRandUser()
			fmt.Printf("[BATTLE REPORT][%v] start pvp battle!\n", enemyU.Account)
			initData.EnemyTeam = enemyU.Team
			initData.EnemyName = fmt.Sprintf("%v", enemyU.Account)
		}
		sendMsg := mainBattle.Report(initData)
		fmt.Println("[BATTLE REPORT]: " + sendMsg)
		mainServer.SendMsg(id, sys, api, sendMsg)
	case MainEvent.CSC_BATTLE_UPDATE_TEAM: //更新隊伍
		agentsStr := strings.Split(msg, "|")
		agents := []model.Agent{}
		for _, v := range agentsStr {
			fmt.Printf("[BATTLE] UPATE PLAYER TEAM :%v \n", v)
			if v == "" {
				continue
			}
			agentStr := strings.Split(v, ",")
			pixel, _ := strconv.Atoi(agentStr[0])
			x, _ := strconv.Atoi(agentStr[1])
			y, _ := strconv.Atoi(agentStr[2])

			xOffset := 2
			agent := model.Agent{
				Pixel: pixel, MapId: 0, Hp: 10, Face: 0,
				Pos: model.Vector2{X: x + xOffset, Y: y},
			}
			agents = append(agents, agent)
		}
		u := userMgr.Users[mainServer.Clients[id].Account]
		userMgr.UpdateTeam(u.Account, agents)
		fmt.Printf("[BATTLE] UPATE PLAYER TEAM FINISH\n")
	case MainEvent.CSC_GACHA: //抽卡
		fmt.Printf("[CSC_GACHA] %v\n", userMgr.UIDToken)
		items := []model.Item{{UID: userMgr.UIDToken, Id: 1000 + rand.Intn(13), Quantity: 1}}
		userMgr.UIDToken++
		sendMsg := ""
		for _, v := range items {
			u := userMgr.Users[mainServer.Clients[id].Account]
			userMgr.AddItem(u.Account, v)
			sendMsg += fmt.Sprintf("%v,%v,%v|", v.UID, v.Id, v.Quantity)
		}
		mainServer.SendMsg(id, sys, MainEvent.CSC_GACHA, sendMsg)
		mainServer.SendMsg(id, sys, MainEvent.SC_BOX_GET_ITEM, sendMsg)
	default:

	}
}
