package redis

/*
import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mediocregopher/radix/v4"

	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/storage"
)

const (
	keyName = "void_players"
)

type playerManager struct {
	Client radix.Client
}

func NewPlayers(client radix.Client) storage.Players {
	return &playerManager{Client: client}
}

func (rm *playerManager) PlayerByID(ctx context.Context, id string) (models.Player, error) {
	var data []byte
	err := rm.Client.Do(ctx, radix.Cmd(&data, "HGET", keyName, id))
	if err != nil {
		log.Printf("err: %s", err)
		return models.Player{}, err
	}

	log.Printf("%s", data)

	var player models.Player
	err = json.Unmarshal(data, &player)
	if err != nil {
		log.Printf("err: %s", err)
		return models.Player{}, err
	}
	log.Printf("%+v", player)

	return player, nil
}

func (rm *playerManager) AllPlayers(ctx context.Context) ([]models.Player, error) {
	var playerList map[string]string

	err := rm.Client.Do(ctx, radix.Cmd(&playerList, "HGETALL", keyName))
	if err != nil {
		return nil, err
	}

	players := make([]models.Player, len(playerList))
	index := 0
	for _, v := range playerList {
		var p models.Player
		err = json.Unmarshal([]byte(v), &p)
		if err != nil {
			return nil, err
		}
		players[index] = p
		index++
	}

	return players, nil
}

func (rm *playerManager) AllPlayersMap(ctx context.Context) (map[string]models.Player, error) {
	allPlayers, err := rm.AllPlayers(ctx)
	if err != nil {
		return nil, err
	}

	players := make(map[string]models.Player)
	for _, player := range allPlayers {
		players[player.ID] = player
	}

	return players, nil
}

func (rm *playerManager) SavePlayer(ctx context.Context, player models.Player) (models.Player, error) {
	player.Updated = time.Now().UTC()
	data, err := json.Marshal(player)
	if err != nil {
		log.Printf("err: %s", err)
		return models.Player{}, err
	}

	err = rm.Client.Do(ctx, radix.Cmd(nil, "HSET", keyName, player.ID, string(data)))
	if err != nil {
		log.Printf("Err: %s", err)
		return models.Player{}, err
	}

	return player, nil
}

func (rm *playerManager) DeletePlayer(ctx context.Context, player models.Player) error {
	if err := rm.Client.Do(ctx, radix.Cmd(nil, "HDEL", keyName, player.ID)); err != nil {
		return err
	}

	return nil
}
*/
