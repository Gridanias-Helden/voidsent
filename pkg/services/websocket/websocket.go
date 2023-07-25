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
		_ = s.CloseWithMsg([]byte("no session found"))
		return
	}

	joinMsg, _ := json.Marshal(wsMessage{
		Body: map[string]string{
			"avatar": session.Avatar,
			"from":   session.Username,
			"room":   "lobby",
			"time":   time.Now().UTC().Format(time.RFC3339),
		},
		Type: "room:join",
	})
	sessionMsg, _ := json.Marshal(wsMessage{
		Body: map[string]string{
			"avatar": session.Avatar,
			"from":   session.Username,
			"time":   time.Now().UTC().Format(time.RFC3339),
		},
		Type: "session",
	})

	s.Set("session", session)
	s.Set("room", "lobby")

	time.Sleep(50 * time.Millisecond)

	_ = s.Write(sessionMsg)
	_ = ws.mel.BroadcastFilter(joinMsg, ws.toRoom("lobby"))
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

	msg, _ := json.Marshal(map[string]any{
		"body": map[string]string{
			"avatar": session.Avatar,
			"from":   session.Username,
			"room":   room,
			"time":   time.Now().UTC().Format(time.RFC3339),
		},
		"type": "room:leave",
	})

	_ = ws.mel.BroadcastFilter(msg, ws.toRoom(room))

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
			newMsg, _ := json.Marshal(map[string]any{
				"body": map[string]string{
					"msg":  chatMsg,
					"from": session.Username,
					"to":   to,
					"time": time.Now().UTC().Format(time.RFC3339),
				},
				"type": "chat:whisper",
			})
			_ = ws.mel.BroadcastFilter(newMsg, ws.toName(to))
			_ = s.Write(newMsg)
		} else {
			newMsg, _ := json.Marshal(map[string]any{
				"body": map[string]string{
					"msg":  chatMsg,
					"from": session.Username,
					"time": time.Now().UTC().Format(time.RFC3339),
				},
				"type": "chat:all",
			})
			_ = ws.mel.BroadcastFilter(newMsg, ws.toRoom(room))
		}
	}
}

func (ws *WebSocket) toRoom(room string) func(*melody.Session) bool {
	return func(s *melody.Session) bool {
		roomStr, ok := s.MustGet("room").(string)
		if !ok {
			return false
		}

		return roomStr == room
	}
}

func (ws *WebSocket) toName(name string) func(*melody.Session) bool {
	return func(s *melody.Session) bool {
		session, ok := s.MustGet("session").(models.Session)
		if !ok {
			return false
		}

		return session.Username == name
	}
}
