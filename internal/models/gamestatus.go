package models

type GameStatus int

const (
	StatusCards GameStatus = iota
	StatusWerewolves
	StatusSacrifice
	StatusDebate
	StatusTribunal
	// StatusThief
	// StatusCupid
	// StatusCouple
	StatusFortuneTeller
	StatustWitch
)
