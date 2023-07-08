package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mediocregopher/radix/v4"

	"github.com/gridanias-helden/voidsent/internal/models"
)

type PlayerManager struct {
	Client radix.Client
}

func NewPlayerManager(client radix.Client) models.PlayerManager {
	return &PlayerManager{Client: client}
}

func (rm *PlayerManager) LoadPlayerByID(ctx context.Context, id string) (*models.Player, error) {
	var data []byte
	err := rm.Client.Do(ctx, radix.Cmd(&data, "HGET", "player", id))
	if err != nil {
		log.Printf("err: %s", err)
		return nil, err
	}

	log.Printf("%s", data)

	var player models.Player
	err = json.Unmarshal(data, &player)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, err
	}
	log.Printf("%+v", player)

	return &player, nil
}

func (rm *PlayerManager) LoadAllPlayers(ctx context.Context) ([]*models.Player, error) {
	return nil, nil
}

func (rm *PlayerManager) SavePlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	data, err := json.Marshal(player)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, err
	}

	err = rm.Client.Do(ctx, radix.Cmd(nil, "HSET", "player", player.ID, string(data)))
	if err != nil {
		log.Printf("Err: %s", err)
		return nil, err
	}

	return player, nil
}

func (rm *PlayerManager) DeletePlayer(ctx context.Context, player *models.Player) error {
	return nil
}
