package yml

import (
	"context"
	"os"

	"github.com/go-yaml/yaml"

	"github.com/christopher-kleine/werewolves/internal/models"
)

type GameManager struct {
	Games []*models.Game `yaml:"games"`
}

func NewGameManager(filename string) (models.GameManager, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var manager *GameManager
	yaml.Unmarshal(data, &manager)

	return manager, nil
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
