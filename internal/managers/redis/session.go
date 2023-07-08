package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mediocregopher/radix/v4"

	"github.com/gridanias-helden/voidsent/internal/models"
)

type SessionManager struct {
	Client radix.Client
}

func NewSessionManager(client radix.Client) models.SessionManager {
	return &SessionManager{Client: client}
}

func (rm *SessionManager) LoadSessionByID(ctx context.Context, id string) (*models.Session, error) {
	var data []byte
	err := rm.Client.Do(ctx, radix.Cmd(&data, "HGET", "session", id))
	if err != nil {
		log.Printf("Sess Err: %v", err)
		return nil, err
	}

	var session models.Session
	err = json.Unmarshal(data, &session)
	if err != nil {
		log.Printf("SessION Err: %v", err)
		return nil, err
	}

	return &session, nil
}

func (rm *SessionManager) LoadAllSessions(ctx context.Context) ([]*models.Session, error) {
	return nil, nil
}

func (rm *SessionManager) SaveSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	err = rm.Client.Do(ctx, radix.Cmd(nil, "HSET", "session", session.ID, string(data)))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (rm *SessionManager) DeleteSession(ctx context.Context, session *models.Session) error {
	return nil
}
