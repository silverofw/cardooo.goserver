package user

import (
	"cardooo/model"
	"fmt"
	"math/rand"
)

type UserMgr struct {
	Users map[int]*model.User

	UIDToken int
}

func InitUserMgr() *UserMgr {
	mgr := UserMgr{
		Users:    make(map[int]*model.User),
		UIDToken: 1000,
	}
	return &mgr
}

func (u *UserMgr) UserLogin(account int, passward string, serverId int) {
	user, ok := u.Users[account]
	if ok {
		fmt.Printf("[USER][UserLogin][%v] Welcome back!\n", account)
		user.ServerId = serverId

	} else {
		fmt.Printf("[USER][UserLogin][%v] Welcome new user!!\n", account)
		newU := model.User{
			Account:  account,
			Passward: passward,
			Items: []model.Item{
				{UID: u.UIDToken, Id: 1, Quantity: 1},
			},
			Team: []model.Agent{
				{MapId: 0, Hp: 10, Face: 0, Pixel: 1, Pos: model.Vector2{X: 4, Y: 2}},
			},
		}
		u.UIDToken++
		u.Users[account] = &newU
	}
}

func (u *UserMgr) UpdateTeam(account int, agents []model.Agent) {
	c := u.Users[account]
	c.Team = agents
	u.Users[account] = c
}

func (u *UserMgr) AddItem(account int, addItem model.Item) {
	c := u.Users[account]
	c.Items = append(c.Items, addItem)
	u.Users[account] = c
}

func (u *UserMgr) MainQuestFinish(account int) {
	c := u.Users[account]
	c.MainQuest++
	u.Users[account] = c
}

func (u *UserMgr) GetRandUser() *model.User {
	l := len(u.Users)
	index := rand.Intn(l)
	i := 0
	for _, v := range u.Users {
		if i == index {
			return v
		}
		i++
	}
	return nil
}
