package models

type GameStatus int

const (
	StatusAwaitingPlayer GameStatus = iota
	StatusCards
	StatusWerewolves
	StatusSacrifice
	StatusDebate
	StatusTribunal
	StatusThief
	StatusCupid
	StatusCouple
	StatusFortuneTeller
	StatustWitch
)
