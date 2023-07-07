package memory

import (
	"context"

	"github.com/gridanias-helden/voidsent/internal/models"
)

type GameManager struct {
	Games []*models.Game `yaml:"games"`
}

func NewGameManager() models.GameManager {
	return &GameManager{
		make([]*models.Game, 0),
	}
}

func (ym *GameManager) LoadGameByID(ctx context.Context, id string) (*models.Game, error) {
	return nil, nil
}

func (ym *GameManager) LoadAllGames(ctx context.Context) ([]*models.Game, error) {
	return nil, nil
}

func (ym *GameManager) SaveGame(ctx context.Context, game *models.Game) (*models.Game, error) {
	return nil, nil
}

func (ym *GameManager) DeleteGame(ctx context.Context, game *models.Game) error {
	return nil
}
