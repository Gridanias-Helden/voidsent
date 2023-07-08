package models

import "context"

type Game struct {
	ID      string          `json:"id"`
	Players map[string]Role `json:"players"`
	Name    string          `json:"name"`
	Status  GameStatus      `json:"status"`
}

type GameManager interface {
	LoadGameByID(ctx context.Context, id string) (*Game, error)
	LoadAllGames(ctx context.Context) ([]*Game, error)
	SaveGame(ctx context.Context, game *Game) (*Game, error)
	DeleteGame(ctx context.Context, game *Game) error
}
