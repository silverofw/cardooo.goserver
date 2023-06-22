package model

type User struct{
	Account int
	Passward string

	ServerId int

	Items []Item
	Team []Agent
}