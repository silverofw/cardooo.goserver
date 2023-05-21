package game

import (
	"fmt"
	"time"
)

type Game struct {
	Uid int
	name string
	frame int
	token int
	
	AgentMap map[int]Agent

	SendMsg func(int, int, int, string)
}

type Agent struct
{
	id int
	pos Pos
	frame int
}

type Pos struct {
	x int
	y int
}

func InitGame() Game{
	g := Game{
		Uid: 7777,
		name: "PixelMonWorld",
		frame: 0,
		token: 0,
		AgentMap: make(map[int]Agent),
	}

	go  g.HandleGame()	

	// test
	g.AddNewAgent(1) 
	g.RemoveAgent(1)

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
		for _,v := range g.AgentMap {
			v.frame++
		}		

		time.Sleep(1000*1000*1000)
	}
}
func (g *Game)AddNewAgent(id int) {		
	_, ok := g.AgentMap[id]
	if ok {
		fmt.Printf("[Game]have same id %v\n",id)
		return
	}

	newAgent := Agent{
		id: id,
		pos: Pos{
			x: 2,
			y: 2,
		},
		frame: 0,
	}
	g.AgentMap[id] = newAgent
	info := fmt.Sprintf("[Game][AddNewAgent] id: %v",id)
	fmt.Println(info)
	if g.SendMsg != nil {
		g.SendMsg(0, 1, 3, info)
		info = fmt.Sprintf("%v,%v,%v",newAgent.id,newAgent.pos.x,newAgent.pos.y)
		g.SendMsg(id, 1, 10, info)
	}
}
func (g *Game)RemoveAgent(id int) {
	fmt.Printf("[Game][RemoveAgent] id: %v, frame: %v\n",id,g.AgentMap[id].frame)
	delete(g.AgentMap, id)
}

