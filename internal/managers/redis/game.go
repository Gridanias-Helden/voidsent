package redis

import (
	"context"
	"encoding/json"

	"github.com/mediocregopher/radix/v4"

	"github.com/gridanias-helden/voidsent/internal/models"
)

type GameManager struct {
	Client radix.Client
}

func NewGameManager(client radix.Client) models.GameManager {
	return &GameManager{Client: client}
}

func (rm *GameManager) LoadGameByID(ctx context.Context, id string) (*models.Game, error) {
	var data []byte
	err := rm.Client.Do(ctx, radix.Cmd(&data, "HGET", "game", id))
	if err != nil {
		return nil, err
	}

	var game models.Game
	err = json.Unmarshal(data, &game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (rm *GameManager) LoadAllGames(ctx context.Context) ([]*models.Game, error) {
	var gameList map[string]string

	err := rm.Client.Do(ctx, radix.Cmd(&gameList, "HGETALL", "game"))
	if err != nil {
		return nil, err
	}

	games := make([]*models.Game, len(gameList))
	index := 0
	for _, v := range gameList {
		var g models.Game
		err = json.Unmarshal([]byte(v), &g)
		if err != nil {
			return nil, err
		}
		games[index] = &g
		index++
	}

	return games, nil
}

func (rm *GameManager) SaveGame(ctx context.Context, game *models.Game) (*models.Game, error) {
	data, err := json.Marshal(game)
	if err != nil {
		return nil, err
	}

	err = rm.Client.Do(ctx, radix.Cmd(nil, "HSET", "game", game.ID, string(data)))
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (rm *GameManager) DeleteGame(ctx context.Context, game *models.Game) error {
	return nil
}
