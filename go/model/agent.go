package model

import (
	Math "cardooo/core"
	MainEvent "cardooo/enum"
)

type Agent struct {
	Id    int
	Hp    int
	Face  int
	Pos   Vector2
	MapId int
	Frame int
	Pixel int
	Team  int
}

func (a *Agent) AddPos(x int, y int) {
	a.Pos.X += x
	a.Pos.Y += y
}

func (a *Agent) Reverse() {
	a.Pos.X = 8 - a.Pos.X
	a.Pos.Y = 8 - a.Pos.Y
}

func (a *Agent) Behit(damage int) {
	a.Hp -= damage
}

func (a Agent) GetOrder(target *Agent) int {
	delX := Math.Abs(target.Pos.X - a.Pos.X)
	delY := Math.Abs(target.Pos.Y - a.Pos.Y)
	if a.Pos.Dis(target.Pos) == 1 {
		if delX > delY {
			if target.Pos.X > a.Pos.X {
				if a.Face != 3 {
					return MainEvent.ORDER_MOVE_RIGHT
				} else {
					return MainEvent.ORDER_ATTACK
				}
			}
			if a.Face != 2 {
				return MainEvent.ORDER_MOVE_LEFT
			} else {
				return MainEvent.ORDER_ATTACK
			}
		} else {
			if target.Pos.Y > a.Pos.Y {
				if a.Face != 0 {
					return MainEvent.ORDER_MOVE_UP
				} else {
					return MainEvent.ORDER_ATTACK
				}
			}
			if a.Face != 1 {
				return MainEvent.ORDER_MOVE_DOWN
			} else {
				return MainEvent.ORDER_ATTACK
			}
		}
	}
	if delX > delY {
		if target.Pos.X > a.Pos.X {
			return MainEvent.ORDER_MOVE_RIGHT
		}
		return MainEvent.ORDER_MOVE_LEFT
	} else {
		if target.Pos.Y > a.Pos.Y {
			return MainEvent.ORDER_MOVE_UP
		}
		return MainEvent.ORDER_MOVE_DOWN
	}
}
