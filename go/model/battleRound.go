package model

type BattleRound struct {
	Orders []Order
}
type Order struct {
	EntityId int
	OrderId  int
	Params   []int
}