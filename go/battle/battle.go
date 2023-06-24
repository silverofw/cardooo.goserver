package battle

import (
	"cardooo/game"
	"cardooo/model"
	"fmt"
)

type Battle struct {
	Uid      int
	reportId int

	battleGame game.Game
}

type BattleInitData struct {
	PlayerTeam []model.Agent
	EnemyTeam  []model.Agent
	EnemyName string
}

func InitBattle() Battle {
	b := Battle{
		Uid:      1,
		reportId: 1,
	}

	return b
}

func (b *Battle) GetGame(initData BattleInitData) game.Game {	

	g := game.InitBattleGame()
	g.Rounds = []model.BattleRound{}
	g.Agents = []model.Agent{}

	tokenId := 10
	for _, a := range initData.PlayerTeam {
		a.Id = tokenId
		tokenId++
		a.Team = 0
		g.AddNewAgent(a)
		g.Agents = append(g.Agents, a)
	}
	for _, a := range initData.EnemyTeam {
		a.Id = tokenId
		tokenId++
		a.Team = 1
		a.Reverse() // switch offset
		g.AddNewAgent(a)
		g.Agents = append(g.Agents, a)
	}

	// 建立動作順位
	index := 0
	actionArray := []int{}
	for _, a := range g.AgentMap {
		actionArray = append(actionArray, a.Id)
	}

	for !g.GameEnd() {
		// select new action agent
		id := actionArray[index]
		index++
		index %= len(actionArray)

		//fmt.Printf("[BATTLE] Round %v Start! id:%v ======\n", len(g.Rounds)+1, id)

		curAgent := g.AgentMap[id]
		if curAgent == nil {
			//fmt.Printf("[BATTLE] curAgent is dead~ id:%v\n", id)
			continue
		}

		target := g.GetTarget(id)
		orderId := 0
		if target != nil {
			orderId = curAgent.GetOrder(target)
		}
		// process order
		g.OnOrder(id, orderId)

		r := model.BattleRound{}
		r.Orders = []model.Order{}
		r.Orders = append(r.Orders, model.Order{
			EntityId: id,
			OrderId:  orderId,
		})

		if target != nil && target.Hp <= 0 {
			r.Orders = append(r.Orders, model.Order{
				EntityId: 1, //戰場控制者
				OrderId:  8102,
				Params:   []int{target.Id},
			})
		}

		g.Rounds = append(g.Rounds, r)

		if len(g.Rounds) > 100 {
			fmt.Printf("[ERROR][BATTLE] round too many!\n")
			break
		}
	}	
	return g
}

func (b *Battle) Report(initData BattleInitData, g game.Game) string {
	sendMsg := ""
	sendMsg += fmt.Sprintf("%v,%v,%v,%v,%s|", 
		b.reportId, g.WinTeam, len(g.Agents), len(g.Rounds), initData.EnemyName)
	for _, a := range g.Agents {
		sendMsg += fmt.Sprintf("%v,%v,%v,%v,%v|", a.Id, a.Pixel, a.Hp, a.Pos.X, a.Pos.Y)
	}

	for _, r := range g.Rounds {
		for _, o := range r.Orders {
			sendMsg += fmt.Sprintf("%v,%v", o.EntityId, o.OrderId)
			if len(o.Params) != 0 {
				for _, p := range o.Params {
					sendMsg += fmt.Sprintf(",%v", p)
				}
			}
			sendMsg += "="
		}
		sendMsg += "|"
	}

	b.reportId++
	return sendMsg
}
