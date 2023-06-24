package model

type User struct{
	Account int
	Passward string

	ServerId int

	MainQuest int

	Items []Item
	Team []Agent
}