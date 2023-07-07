package services

import (
	"github.com/gridanias-helden/voidsent/internal/managers/memory"
	"github.com/gridanias-helden/voidsent/internal/models"
)

type Memory struct {
	models.PlayerManager
	models.GameManager
}

func NewMemory() Service {
	return &YAML{
		PlayerManager: memory.NewPlayerManager(),
		GameManager:   memory.NewGameManager(),
	}
}
