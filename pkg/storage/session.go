package storage

import (
	"context"

	"github.com/gridanias-helden/voidsent/pkg/models"
)

type Sessions interface {
	SessionByID(ctx context.Context, id string) (models.Session, error)
	SaveSession(ctx context.Context, session models.Session) (models.Session, error)
	DeleteSession(ctx context.Context, session models.Session) error
}
