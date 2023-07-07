package models

import "context"

type Player struct {
	ID     string
	Name   string
	Avatar string
}

type PlayerManager interface {
	LoadPlayerByID(ctx context.Context, id string) (*Player, error)
	LoadAllPlayers(ctx context.Context) ([]*Player, error)
	SavePlayer(ctx context.Context, player *Player) (*Player, error)
	DeletePlayer(ctx context.Context, player *Player) error
}
