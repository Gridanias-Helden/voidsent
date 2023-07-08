package services

import (
	"github.com/mediocregopher/radix/v4"

	"github.com/gridanias-helden/voidsent/internal/managers/redis"
	"github.com/gridanias-helden/voidsent/internal/models"
)

type Redis struct {
	models.PlayerManager
	models.GameManager
	models.SessionManager
}

func NewRedis(client radix.Client) Service {
	return &YAML{
		PlayerManager:  redis.NewPlayerManager(client),
		GameManager:    redis.NewGameManager(client),
		SessionManager: redis.NewSessionManager(client),
	}
}
