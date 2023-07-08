package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oklog/ulid"
	"github.com/olahol/melody"

	"github.com/gridanias-helden/voidsent/internal/middleware"
	"github.com/gridanias-helden/voidsent/internal/models"
	"github.com/gridanias-helden/voidsent/internal/services"
)

type newGameMessage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WebSocketMessage struct {
	Type    string          `json:"type"`
	NewGame *newGameMessage `json:"newGame"`
	Lobby   []*models.Game  `json:"lobby"`
}

type WebSocket struct {
	Melody  *melody.Melody
	Service services.Service
}

func (ws *WebSocket) HTTPRequest(w http.ResponseWriter, r *http.Request) {
	ws.Melody.HandleRequest(w, r)
}

func (ws *WebSocket) Connect(s *melody.Session) {
	session, ok := s.Request.Context().Value(middleware.SessionKey).(*models.Session)
	if !ok {
		s.Close()
		return
	}

	allGames, err := ws.Service.LoadAllGames(s.Request.Context())
	if err != nil {
		log.Printf("load lobby: %v", err)
		return
	}

	// Send lobby
	msg, err := json.Marshal(WebSocketMessage{
		Type:  "lobby",
		Lobby: allGames,
	})

	if err != nil {
		log.Printf("encode err %s", err)
		return
	}

	log.Printf("lobby: %s", string(msg))

	time.Sleep(100 * time.Millisecond)
	s.Set("session", session)
	if err := s.Write(msg); err != nil {
		log.Printf("write err %s", err)
		return
	}

	log.Printf("connected: %s", s.Request.RemoteAddr)
}

func (ws *WebSocket) Message(s *melody.Session, msg []byte) {
	clientIntf, ok := s.Get("session")
	if !ok {
		log.Printf("ws message: no session preset")
		return
	}

	session, ok := clientIntf.(*models.Session)
	if !ok {
		log.Printf("ws message: session invalid")
		return
	}

	var message WebSocketMessage
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Printf("decode err %s", err)
		return
	}

	switch message.Type {
	case "newGame":
		game := &models.Game{
			ID:     ulid.MustNew(uint64(time.Now().UnixMilli()), rand.Reader).String(),
			Status: models.StatusAwaitingPlayer,
			Name:   message.NewGame.Name,
			Players: map[string]models.Role{
				session.ID: models.RoleUndecided,
			},
		}

		game, err = ws.Service.SaveGame(s.Request.Context(), game)
		if err != nil {
			log.Printf("new game: %v", err)
			return
		}

		session.Game = game.ID
		session, err = ws.Service.SaveSession(s.Request.Context(), session)
		if err != nil {
			log.Printf("update session: %v", err)
			return
		}

		err = ws.Melody.BroadcastFilter([]byte("New GAME CREATED"), func(s *melody.Session) bool {
			clientIntf, ok := s.Get("session")
			if !ok {
				log.Printf("ws message: no session preset")
				return false
			}

			session, ok := clientIntf.(*models.Session)
			if !ok {
				log.Printf("ws message: session invalid")
				return false
			}

			return session.Game == ""
		})
		if err != nil {
			log.Printf("ws broadcast %s", err)
		}

		content := fmt.Sprintf(`{"type": "joinRoom", "id": %s }`, game.ID)
		s.Write([]byte(content))
	}

	log.Printf("received (text): %s", string(msg))
}

func (ws *WebSocket) MessageBinary(s *melody.Session, msg []byte) {
	log.Printf("received (binary): %s", string(msg))
}
