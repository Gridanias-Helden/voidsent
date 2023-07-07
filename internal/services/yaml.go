package services

import (
	"github.com/gridanias-helden/voidsent/internal/managers/yml"
	"github.com/gridanias-helden/voidsent/internal/models"
)

type YAML struct {
	models.PlayerManager
	models.GameManager
	models.SessionManager
}

func NewYAML(filename string) (Service, error) {
	playerManager, err := yml.NewPlayerManager(filename)
	if err != nil {
		return nil, err
	}

	gameManager, err := yml.NewGameManager(filename)
	if err != nil {
		return nil, err
	}

	return &YAML{
		PlayerManager: playerManager,
		GameManager:   gameManager,
		//SessionManager: sessionManager,
	}, nil
}
