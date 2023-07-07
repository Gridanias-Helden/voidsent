package services

import "github.com/gridanias-helden/voidsent/internal/models"

type Service interface {
	models.PlayerManager
	models.GameManager
}
