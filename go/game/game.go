package game

import (
	MainEvent "cardooo/enum"
	"cardooo/model"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Game struct {
	Uid        int
	name       string
	frame      int
	token      int
	autoCd     int
	agentToken int
	StartMapId int
	MapId      int

	Agents []model.Agent
	Rounds []model.BattleRound

	WinTeam int

	AgentMap map[int]*model.Agent

	SendMsg          func(int, int, int, string)
	SendBroadcastMsg func(int, int, int, string)
}

var tils = [32]model.Vector2{
	{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 0}, {X: 7, Y: 0}, {X: 8, Y: 0},
	{X: 0, Y: 8}, {X: 1, Y: 8}, {X: 2, Y: 8}, {X: 3, Y: 8}, {X: 4, Y: 8}, {X: 5, Y: 8}, {X: 6, Y: 8}, {X: 7, Y: 8}, {X: 8, Y: 8},
	{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}, {X: 0, Y: 7},
	{X: 8, Y: 1}, {X: 8, Y: 2}, {X: 8, Y: 3}, {X: 8, Y: 4}, {X: 8, Y: 5}, {X: 8, Y: 6}, {X: 8, Y: 7},
}
var facePos = [4]model.Vector2{
	{X: 0, Y: 1}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: 1, Y: 0},
}

func InitBattleGame() Game {
	g := Game{
		Uid:      7777,
		name:     "PixelMonWorld",
		frame:    0,
		token:    0,
		autoCd:   3,
		MapId:    7007,
		AgentMap: make(map[int]*model.Agent),
		WinTeam:  -1,
	}
	return g
}

func InitLobyGame(SendMsg func(int, int, int, string), SendBroadcastMsg func(int, int, int, string)) Game {
	g := Game{
		Uid:              7777,
		name:             "PixelMonWorld",
		frame:            0,
		token:            0,
		autoCd:           3,
		MapId:            7007,
		AgentMap:         make(map[int]*model.Agent),
		SendMsg:          SendMsg,
		SendBroadcastMsg: SendBroadcastMsg,
	}
	go g.LobyMode()
	return g
}

func (g *Game) LobyMode() {
	dt := time.Now()
	fmt.Println("[Game] Start!", dt.String())

	for {
		g.frame = g.frame + 1
		var hour int = g.frame / 60 / 60
		var minute int = g.frame / 60 % 60
		var secend int = g.frame % 60
		if secend == 0 {
			log.Printf("[Game][%v][%v:%v:%v]\n", g.frame, hour, minute, secend)
		}

		g.AutoAgentCreator()
		for _, v := range g.AgentMap {
			v.Frame++
		}

		// 等待一秒
		time.Sleep(time.Second)
	}
}
func (g *Game) AutoAgentCreator() {
	g.autoCd--
	if g.autoCd > 0 {
		return
	}

	g.autoCd = 10
	if len(g.AgentMap) < 3 {
		fmt.Printf("[Game][AutoAgentCreator][%v] auto create new agent! \n", len(g.AgentMap))
		g.agentToken++
		a := g.AddNewAgentById(g.agentToken)
		info := fmt.Sprintf("%v,%v,%v", a.Id, a.Pos.X, a.Pos.Y)
		g.SendBroadcastMsg(0, 1, 1101, info)
	}
}

func (g *Game) AddNewAgentById(id int) *model.Agent {
	_, ok := g.AgentMap[id]
	if ok {
		fmt.Printf("[Game]have same id %v\n", id)
		return g.AgentMap[id]
	}

	newAgent := &model.Agent{
		Id:    id,
		MapId: g.MapId,
		Hp:    10,
		Face:  0,
		Pos: model.Vector2{
			X: rand.Intn(7) + 1,
			Y: rand.Intn(7) + 1,
		},
		Frame: 0,
	}
	g.AgentMap[id] = newAgent

	fmt.Printf("[Game][AddNewAgent] id: %v\n", id)
	return g.AgentMap[id]
}
func (g *Game) AddNewAgent(agent model.Agent) *model.Agent {
	_, ok := g.AgentMap[agent.Id]
	if ok {
		fmt.Printf("[Game]have same id %v\n", agent.Id)
		return g.AgentMap[agent.Id]
	}

	g.AgentMap[agent.Id] = &agent

	fmt.Printf("[Game][AddNewAgent] id: %v\n", agent.Id)
	return &agent
}

func (g *Game) RemoveAgent(id int) {
	fmt.Printf("[Game][RemoveAgent] id: %v, frame: %v\n", id, g.AgentMap[id].Frame)
	delete(g.AgentMap, id)
}

func (g *Game) GetTarget(id int) *model.Agent {
	curAgent := g.AgentMap[id]
	if curAgent == nil {
		return nil
	}
	return g.getNearestAgent(curAgent)
}

func (g *Game) OnOrder(id int, order int) {
	fmt.Printf("[Game][OnOrder] id: %v, order: %v\n", id, order)
	switch order {
	case MainEvent.ORDER_MOVE_UP:
		g.AgentMap[id].Face = 0
		var checkPos = g.AgentMap[id].Pos.Add(facePos[0])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(0, 1)
			fmt.Printf("[Game][OnOrder][8001] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case MainEvent.ORDER_MOVE_DOWN:
		g.AgentMap[id].Face = 1
		var checkPos = g.AgentMap[id].Pos.Add(facePos[1])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(0, -1)
			fmt.Printf("[Game][OnOrder][8002] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case MainEvent.ORDER_MOVE_LEFT:
		g.AgentMap[id].Face = 2
		var checkPos = g.AgentMap[id].Pos.Add(facePos[2])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(-1, 0)
			fmt.Printf("[Game][OnOrder][8003] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case MainEvent.ORDER_MOVE_RIGHT:
		g.AgentMap[id].Face = 3
		var checkPos = g.AgentMap[id].Pos.Add(facePos[3])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(1, 0)
			fmt.Printf("[Game][OnOrder][8004] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case 8009: // attack
		a := g.getAgent(facePos[g.AgentMap[id].Face].Add(g.AgentMap[id].Pos))
		if a == nil {
			fmt.Printf("[Game][OnOrder][8009] face is nil \n")
			return
		}
		a.Behit(1)
		fmt.Printf("[Game][OnOrder][8009]id:%v ,Hp: %v\n", a.Id, a.Hp)
		if a.Hp <= 0 {
			sendMsg := fmt.Sprintf("%v,%v,%v|", a.Id, a.Pos.X, a.Pos.Y)
			if g.SendBroadcastMsg != nil {
				g.SendBroadcastMsg(0, 1, 1102, sendMsg)
			}
			g.RemoveAgent(a.Id)
		}
	}
}

func (g *Game) getAgent(pos model.Vector2) *model.Agent {
	for _, v := range g.AgentMap {
		if v.Pos.X == pos.X && v.Pos.Y == pos.Y {
			return v
		}
	}
	return nil
}

func (g *Game) getNearestAgent(agent *model.Agent) *model.Agent {
	var target *model.Agent
	minDis := 99999
	for _, v := range g.AgentMap {
		if v.Team != agent.Team && v.Id != agent.Id {
			dis := agent.Pos.Dis(v.Pos)
			if minDis > dis {
				minDis = dis
				target = v
			}
		}
	}
	return target
}

func CheckCanPass(pos model.Vector2) bool {
	for _, v := range tils {
		if v.X == pos.X && v.Y == pos.Y {
			return false
		}
	}
	return true
}

func (g *Game) GameEnd() bool {
	team := -1

	for _, a := range g.AgentMap {
		if team == -1 {
			team = a.Team
			continue
		}
		if team != a.Team {
			return false
		}
	}
	//決定勝利隊伍
	g.WinTeam = team

	return true
}
