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

type BattleRound struct {
	Orders []Order
}
type Order struct {
	EntityId int
	OrderId  int
	Params   []int
}

func InitBattle() Battle {
	b := Battle{
		Uid:      1,
		reportId: 1,
	}

	return b
}

func (b *Battle) Report(initData BattleInitData) string {
	agents := []model.Agent{}
	rounds := []BattleRound{}

	b.battleGame = game.InitBattleGame()

	tokenId := 10
	for _, a := range initData.PlayerTeam {
		a.Id = tokenId
		tokenId++
		a.Team = 1
		b.battleGame.AddNewAgent(a)
		agents = append(agents, a)
	}
	for _, a := range initData.EnemyTeam {
		a.Id = tokenId
		tokenId++
		a.Team = 2
		a.Reverse() // switch offset
		b.battleGame.AddNewAgent(a)
		agents = append(agents, a)
	}

	// 建立動作順位
	index := 0
	actionArray := []int{}
	for _, a := range b.battleGame.AgentMap {
		actionArray = append(actionArray, a.Id)
	}

	for !b.battleGame.GameEnd() {
		// select new action agent
		id := actionArray[index]
		index++
		index %= len(actionArray)

		fmt.Printf("[BATTLE] Round %v Start! id:%v ======\n", len(rounds)+1, id)

		curAgent := b.battleGame.AgentMap[id]
		if curAgent == nil {
			fmt.Printf("[BATTLE] curAgent is dead~ id:%v\n", id)
			continue
		}

		target := b.battleGame.GetTarget(id)
		orderId := 0
		if target != nil {
			orderId = curAgent.GetOrder(target)
		}
		// process order
		b.battleGame.OnOrder(id, orderId)

		r := BattleRound{}
		r.Orders = []Order{}
		r.Orders = append(r.Orders, Order{
			EntityId: id,
			OrderId:  orderId,
		})

		if target != nil && target.Hp <= 0 {
			r.Orders = append(r.Orders, Order{
				EntityId: 1, //戰場控制者
				OrderId:  8102,
				Params:   []int{target.Id},
			})
		}

		rounds = append(rounds, r)

		if len(rounds) > 100 {
			fmt.Printf("[ERROR][BATTLE] round too many!\n")
			break
		}
	}

	sendMsg := ""
	sendMsg += fmt.Sprintf("%v,%v,%v,%s|", b.reportId, len(agents), len(rounds), initData.EnemyName)
	for _, a := range agents {
		sendMsg += fmt.Sprintf("%v,%v,%v,%v,%v|", a.Id, a.Pixel, a.Hp, a.Pos.X, a.Pos.Y)
	}

	for _, r := range rounds {
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
