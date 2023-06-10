package model

import (
	Math "cardooo/core"
)

type Vector2 struct {
	X int
	Y int
}
type Vector3 struct {
	X int
	Y int
	Z int
}

func (p Vector2) Add(pos Vector2) Vector2 {
	p.X += pos.X
	p.Y += pos.Y
	return p
}

func (p Vector2) Dis(target Vector2) int {
	dis := 0
	dis += Math.Abs(target.X - p.X)
	dis += Math.Abs(target.Y - p.Y)
	return dis
}
