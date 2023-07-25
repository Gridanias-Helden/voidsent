package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/olahol/melody"

	"github.com/gridanias-helden/voidsent/pkg/middleware"
	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/services"
	"github.com/gridanias-helden/voidsent/pkg/storage"
)

type WebSocket struct {
	sessions storage.Sessions
	broker   *services.Broker
	mel      *melody.Melody
}

type wsMessage struct {
	Type string            `json:"type"`
	Body map[string]string `json:"body"`
}

func New(sessions storage.Sessions, broker *services.Broker, mel *melody.Melody) *WebSocket {
	ws := &WebSocket{
		sessions: sessions,
		broker:   broker,
		mel:      mel,
	}

	mel.HandleConnect(ws.Connect)
	mel.HandleDisconnect(ws.Disconnect)
	mel.HandleMessage(ws.Message)

	return ws
}

func (ws *WebSocket) HTTPRequest(w http.ResponseWriter, r *http.Request) {
	_ = ws.mel.HandleRequest(w, r)
}

func (ws *WebSocket) Connect(s *melody.Session) {
	session, ok := s.Request.Context().Value(middleware.SessionKey).(models.Session)
	if !ok {
		s.CloseWithMsg([]byte("no session found"))
		return
	}

	joinMsg := wsMessage{
		Body: map[string]string{
			"avatar": session.Avatar,
			"name":   session.Username,
			"room":   "lobby",
			"time":   time.Now().UTC().Format(time.RFC3339),
		},
		Type: "room:join",
	}
	sessionMsg, _ := json.Marshal(wsMessage{
		Body: map[string]string{
			"avatar": session.Avatar,
			"name":   session.Username,
			"time":   time.Now().UTC().Format(time.RFC3339),
		},
		Type: "session",
	})

	s.Set("session", session)
	s.Set("room", "lobby")

	time.Sleep(50 * time.Millisecond)

	s.Write(sessionMsg)
	ws.BroadcastRoom(joinMsg, "lobby")
}

func (ws *WebSocket) Disconnect(s *melody.Session) {
	session, ok := s.MustGet("session").(models.Session)
	if !ok {
		return
	}

	room, ok := s.MustGet("room").(string)
	if !ok {
		return
	}

	msg := map[string]any{
		"body": map[string]string{
			"avatar": session.Avatar,
			"name":   session.Username,
			"room":   room,
			"time":   time.Now().UTC().Format(time.RFC3339),
		},
		"type": "room:leave",
	}

	ws.BroadcastRoom(msg, room)

	log.Printf("%s left", session.Username)
}

func (ws *WebSocket) Message(s *melody.Session, msg []byte) {
	session, ok := s.MustGet("session").(models.Session)
	if !ok {
		return
	}

	room, ok := s.MustGet("room").(string)
	if !ok {
		return
	}

	var message wsMessage
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Printf("got message from %s (%s): %s with error: %v", session.ID, session.Username, msg, err)
		return
	}

	if room != "lobby" {
		ws.broker.Send(session.ID, room, "msg", message.Body)
		return
	}

	switch message.Type {
	case "chat":
		chatMsg, ok := message.Body["msg"]
		if !ok {
			return
		}

		if to, ok := message.Body["to"]; ok {
			// Whisper to specific user
			log.Printf("/w to %s", to)
		} else {
			ws.BroadcastRoom([]byte(chatMsg), room)
		}
	}
}

func (ws *WebSocket) BroadcastRoom(msg any, room string) {
	msgBytes, _ := json.Marshal(msg)

	log.Printf("Msg: %s, room: %s", msg, room)
	ws.mel.BroadcastFilter(msgBytes, func(s *melody.Session) bool {
		roomStr, ok := s.MustGet("room").(string)
		if !ok {
			return false
		}

		return roomStr == room
	})
}
