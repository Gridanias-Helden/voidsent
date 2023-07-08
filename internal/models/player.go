package models

import "context"

type Player struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type PlayerManager interface {
	LoadPlayerByID(ctx context.Context, id string) (*Player, error)
	LoadAllPlayers(ctx context.Context) ([]*Player, error)
	SavePlayer(ctx context.Context, player *Player) (*Player, error)
	DeletePlayer(ctx context.Context, player *Player) error
}
