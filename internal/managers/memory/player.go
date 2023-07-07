package memory

import (
	"context"
	"fmt"

	"github.com/gridanias-helden/voidsent/internal/models"
)

type PlayerManager struct {
	Players map[string]*models.Player `yaml:"players"`
}

func NewPlayerManager() models.PlayerManager {
	return &PlayerManager{
		Players: make(map[string]*models.Player, 0),
	}
}

func (mm *PlayerManager) LoadPlayerByID(ctx context.Context, id string) (*models.Player, error) {
	player, ok := mm.Players[id]
	if !ok {
		return nil, fmt.Errorf("player %q not found", player)
	}

	return nil, nil
}

func (mm *PlayerManager) LoadAllPlayers(ctx context.Context) ([]*models.Player, error) {
	var players []*models.Player
	for _, player := range mm.Players {
		players = append(players, player)
	}

	return players, nil
}

func (mm *PlayerManager) SavePlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	mm.Players[player.ID] = player

	return player, nil
}

func (mm *PlayerManager) DeletePlayer(ctx context.Context, player *models.Player) error {
	delete(mm.Players, player.ID)
	return nil
}
