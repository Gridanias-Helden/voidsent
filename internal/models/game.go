package models

import "context"

type Game struct {
	Players map[string]Role
	Status  GameStatus
}

type GameManager interface {
	LoadGameByID(ctx context.Context, id string) (*Game, error)
	LoadAllGames(ctx context.Context) ([]*Game, error)
	SaveGame(ctx context.Context, game *Game) (*Game, error)
	DeleteGame(ctx context.Context, game *Game) error
}
