package handlers

import (
	"log"
	"net/http"

	"github.com/olahol/melody"

	"github.com/gridanias-helden/voidsent/internal/middleware"
	"github.com/gridanias-helden/voidsent/internal/models"
)

type WebSocketMessage struct {
	Type string `json:"type"`
}

type WebSocket struct {
	Melody *melody.Melody
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

	s.Set("session", session)
	log.Printf("connected: %s", s.Request.RemoteAddr)
}

func (ws *WebSocket) Message(s *melody.Session, msg []byte) {
	log.Printf("received (text): %s", string(msg))
}

func (ws *WebSocket) MessageBinary(s *melody.Session, msg []byte) {
	log.Printf("received (binary): %s", string(msg))
}
