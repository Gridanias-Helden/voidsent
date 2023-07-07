package yml

import (
	"context"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/gridanias-helden/voidsent/internal/models"
)

type PlayerManager struct {
	Players []*models.Player `yaml:"players"`
}

func NewPlayerManager(filename string) (models.PlayerManager, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var manager *PlayerManager
	yaml.Unmarshal(data, &manager)

	return manager, nil
}

func (ym *PlayerManager) LoadPlayerByID(ctx context.Context, id string) (*models.Player, error) {
	return nil, nil
}

func (ym *PlayerManager) LoadAllPlayers(ctx context.Context) ([]*models.Player, error) {
	return nil, nil
}

func (ym *PlayerManager) SavePlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	return nil, nil
}

func (ym *PlayerManager) DeletePlayer(ctx context.Context, player *models.Player) error {
	return nil
}
