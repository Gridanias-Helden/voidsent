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
	var playerList map[string]string

	err := rm.Client.Do(ctx, radix.Cmd(&playerList, "HGETALL", "player"))
	if err != nil {
		return nil, err
	}

	players := make([]*models.Player, len(playerList))
	index := 0
	for _, v := range playerList {
		var p models.Player
		err = json.Unmarshal([]byte(v), &p)
		if err != nil {
			return nil, err
		}
		players[index] = &p
		index++
	}

	return players, nil
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
	if err := rm.Client.Do(ctx, radix.Cmd(nil, "HDEL", "player", player.ID)); err != nil {
		return err
	}

	return nil
}
