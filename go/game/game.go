package game

import (
	"fmt"
	"time"
	"math/rand"
)

type Game struct {
	Uid int
	name string
	frame int
	token int
	autoCd int
	agentToken int
	StartMapId int
	MapId int
	
	AgentMap map[int]*Agent

	SendMsg func(int, int, int, string)
	SendBroadcastMsg func(int, int, int, string)
}

type Agent struct
{
	Id int
	Hp int
	Face int
	Pos Pos
	MapId int
	Frame int	
}

type Pos struct {
	X int
	Y int
}

var tils = [32]Pos { 
	Pos{ X: 0, Y: 0,}, Pos{ X: 1, Y: 0,}, Pos{ X: 2, Y: 0,}, Pos{ X: 3, Y: 0,}, Pos{ X: 4, Y: 0,}, Pos{ X: 5, Y: 0,}, Pos{ X: 6, Y: 0,}, Pos{ X: 7, Y: 0,}, Pos{ X: 8, Y: 0,}, 
	Pos{ X: 0, Y: 8,}, Pos{ X: 1, Y: 8,}, Pos{ X: 2, Y: 8,}, Pos{ X: 3, Y: 8,}, Pos{ X: 4, Y: 8,}, Pos{ X: 5, Y: 8,}, Pos{ X: 6, Y: 8,}, Pos{ X: 7, Y: 8,}, Pos{ X: 8, Y: 8,}, 
	Pos{ X: 0, Y: 1,}, Pos{ X: 0, Y: 2,}, Pos{ X: 0, Y: 3,}, Pos{ X: 0, Y: 4,}, Pos{ X: 0, Y: 5,}, Pos{ X: 0, Y: 6,}, Pos{ X: 0, Y: 7,},
	Pos{ X: 8, Y: 1,}, Pos{ X: 8, Y: 2,}, Pos{ X: 8, Y: 3,}, Pos{ X: 8, Y: 4,}, Pos{ X: 8, Y: 5,}, Pos{ X: 8, Y: 6,}, Pos{ X: 8, Y: 7,},  	
}
var facePos = [4]Pos { 
	Pos{ X: 0, Y: 1,}, Pos{ X: 0, Y: -1,}, Pos{ X: -1, Y: 0,}, Pos{ X: 1, Y: 0,}, 
}

func (a *Agent)AddPos(x int, y int){
	a.Pos.X += x
	a.Pos.Y += y
}

func (a *Agent)Behit(damage int){
	a.Hp -= damage
}

func InitGame(SendMsg func(int, int, int, string),SendBroadcastMsg func(int, int, int, string)) Game{
	g := Game{
		Uid: 7777,
		name: "PixelMonWorld",
		frame: 0,
		token: 0,
		autoCd: 3,
		MapId: 7007,
		AgentMap: make(map[int]*Agent),
		SendMsg: SendMsg,
		SendBroadcastMsg: SendBroadcastMsg,
	}

	go  g.HandleGame()

	return g
}

func (g *Game)HandleGame() {
	dt := time.Now()
	fmt.Println("[Game] Start!", dt.String())

	for {
		g.frame = g.frame + 1
		var hour int = g.frame / 60 / 60
		var minute int = g.frame / 60 % 60
		var secend int = g.frame % 60 % 60
		if secend == 0 {
			dt = time.Now()
			fmt.Printf("[Game][%v][%v:%v:%v]%s\n", g.frame, hour, minute, secend, dt.String())
		}

		g.AutoAgentCreator()
		for _,v := range g.AgentMap {
			v.Frame++
		}		

		time.Sleep(1000*1000*1000)
	}
}
func (g *Game)AutoAgentCreator() {
	g.autoCd--
	if g.autoCd > 0 {
		return
	}

	g.autoCd = 10
	if len(g.AgentMap) < 6 {
		fmt.Printf("[Game][AutoAgentCreator][%v] auto create new agent! \n", len(g.AgentMap))
		g.agentToken++;
		a := g.AddNewAgent(g.agentToken)		
		info := fmt.Sprintf("%v,%v,%v",a.Id, a.Pos.X, a.Pos.Y)
		g.SendBroadcastMsg(0, 1, 1101, info)
	}
}

func (g *Game)AddNewAgent(id int) *Agent{		
	_, ok := g.AgentMap[id]
	if ok {
		fmt.Printf("[Game]have same id %v\n",id)
		return g.AgentMap[id]
	}

	newAgent := &Agent{
		Id: id,
		MapId: g.MapId,
		Hp: 10,
		Face: 0,
		Pos: Pos{
			X: rand.Intn(7) + 1,
			Y: rand.Intn(7) + 1,
		},
		Frame: 0,
	}
	g.AgentMap[id] = newAgent
	
	info := fmt.Sprintf("[Game][AddNewAgent] id: %v",id)
	fmt.Println(info)
	return g.AgentMap[id]
}
func (g *Game)RemoveAgent(id int) {
	fmt.Printf("[Game][RemoveAgent] id: %v, frame: %v\n", id, g.AgentMap[id].Frame)
	delete(g.AgentMap, id)
}
func (g *Game)OnOrder(id int, order int) {
	fmt.Printf("[Game][OnOrder] id: %v, order: %v\n", id, order)
	switch order{
	case 8001:
		g.AgentMap[id].Face = 0
		var checkPos = g.AgentMap[id].Pos.Add(facePos[0])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(0,1)
			fmt.Printf("[Game][OnOrder][8001] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case 8002:
		g.AgentMap[id].Face = 1
		var checkPos = g.AgentMap[id].Pos.Add(facePos[1])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(0,-1)
			fmt.Printf("[Game][OnOrder][8002] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case 8003:
		g.AgentMap[id].Face = 2
		var checkPos = g.AgentMap[id].Pos.Add(facePos[2])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(-1,0)
			fmt.Printf("[Game][OnOrder][8003] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case 8004:
		g.AgentMap[id].Face = 3
		var checkPos = g.AgentMap[id].Pos.Add(facePos[3])
		if CheckCanPass(checkPos) && g.getAgent(checkPos) == nil {
			g.AgentMap[id].AddPos(1,0)
			fmt.Printf("[Game][OnOrder][8004] x: %v, y: %v\n", g.AgentMap[id].Pos.X, g.AgentMap[id].Pos.Y)
		}
	case 8009: // attack		
		a := g.getAgent(facePos[g.AgentMap[id].Face].Add(g.AgentMap[id].Pos))
		if a == nil {
			fmt.Printf("[Game][OnOrder][8009] face is nil \n")
			return
		}
		a.Behit(1)
		fmt.Printf("[Game][OnOrder][8009] Hp: %v\n", a.Hp)
		if a.Hp <= 0 {
			sendMsg := fmt.Sprintf("%v,%v,%v|", a.Id, a.Pos.X, a.Pos.Y)	
			g.SendBroadcastMsg(0, 1, 1102, sendMsg)
			g.RemoveAgent(a.Id)
		}
	}
}

func (g *Game)getAgent(pos Pos) *Agent {
	for _,v := range g.AgentMap {
		if(v.Pos.X == pos.X && v.Pos.Y == pos.Y) {
			return v
		}
	}	
	return nil
}

func (p Pos)Add(pos Pos) Pos {
	p.X += pos.X
	p.Y += pos.Y
	return p
}

func CheckCanPass(pos Pos) bool {
	for _, v := range tils {
		if v.X == pos.X && v.Y == pos.Y {
			return false
		}
	}
	return true
}