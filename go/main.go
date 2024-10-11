package main

import (
	"cardooo/battle"
	"cardooo/battleRoomMgr"
	"cardooo/common"
	MainEvent "cardooo/enum"
	"cardooo/game"
	"cardooo/model"
	"cardooo/server"
	"cardooo/user"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

var userMgr *user.UserMgr
var mainServer *server.Server
var mainGame game.Game
var mainBattle battle.Battle
var mainBattleRoomMgr *battleRoomMgr.BattleRoomMgr

func main() {
	fmt.Println("[BMC][1.00] Server Start...")
	// 初始化 User Manager
	userMgr = user.InitUserMgr()
	mainServer = server.InitServer(AddNewAgent, RemoveAgent, ClientCommand)
	mainGame = game.InitLobyGame(mainServer.SendMsg, mainServer.BroadcastMessage)
	mainBattleRoomMgr = battleRoomMgr.NewBattleRoomMgr()

	mainServer.StartTCP()
	mainBattle = battle.InitBattle()

	log.Println("[BMC] Server Started Successfully.")
}

func AddNewAgent(id int) {
	mainGame.AddNewAgentById(id)
	ClientCommand(id, 1, MainEvent.CSC_SERVER_STATE, "")
}

func RemoveAgent(c *common.Client) {
	ServerCommand(c.Id, 1, 1002, "")
	mainGame.RemoveAgent(c.Id)
	mainBattleRoomMgr.LeaveRoom(c)
}

func ServerCommand(id int, sys int, api int, msg string) {
	fmt.Printf("[BMC][ServerCommand] %v,%v,%v,%s\n", id, sys, api, msg)

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
	fmt.Printf("[BMC][ClientCommand] %v,%v,%v,%s\n", id, sys, api, msg)

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
		u := userMgr.Users[mainServer.Clients[id].Account]
		sendMsg := fmt.Sprintf("%v,%v,%v,%v|", id, u.MainQuest, v.Pos.X, v.Pos.Y)
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
				{MapId: 0, Hp: 2 + u.MainQuest*10, Face: 0, Pixel: 2001, Pos: model.Vector2{X: 1, Y: 2}},
				{MapId: 0, Hp: 2 + u.MainQuest*10, Face: 0, Pixel: 2001, Pos: model.Vector2{X: 4, Y: 1}},
			}
			initData.EnemyName = fmt.Sprintf("MAIN %v", u.MainQuest)
		case 1:
			enemyU := userMgr.GetRandUser()
			fmt.Printf("[BATTLE REPORT][%v] start pvp battle!\n", enemyU.Account)
			initData.EnemyTeam = enemyU.Team
			initData.EnemyName = fmt.Sprintf("%v", enemyU.Account)
		}
		g := mainBattle.GetGame(initData)
		fmt.Printf("[BATTLE REPORT] team %v win!\n", g.WinTeam)
		//『ＴＯＤＯ』主線推進邏輯
		if g.WinTeam == 0 {
			u.MainQuest++
			userMgr.MainQuestFinish(mainServer.Clients[id].Account)
			fmt.Printf("[MAIN QQUEST] main quest %v clear!\n", u.MainQuest)
		}
		sendMsg := mainBattle.Report(initData, g)
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
	case MainEvent.CSC_SUMMON_CHESS:
		fmt.Printf("[CSC_SUMMON_CHESS] %v\n", msg)
		targetTick := mainBattleRoomMgr.Rooms[1000].Tick + 60
		mainServer.BroadcastMessage(-1, sys, api, fmt.Sprintf("%v,%v", msg, targetTick))
	case MainEvent.CSC_ENTER_ROOM:
		mainBattleRoomMgr.EnterRoom(mainServer.Clients[id])
	case MainEvent.CSC_LEAVE_ROOM:
		mainBattleRoomMgr.LeaveRoom(mainServer.Clients[id])
	default:

	}
}
