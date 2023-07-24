package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gridanias-helden/voidsent/pkg/models"
)

type newVoidGameMessage struct {
	Name  string          `json:"name"`
	Roles map[string]bool `json:"roles"`
}

type IncomingWebSocketMessage struct {
	Type     string              `json:"type"`
	Voidsent *newVoidGameMessage `json:"voidsent,omitempty"`
}

type WSConn struct {
	ID      string
	Conn    *websocket.Conn
	Msg     chan []byte
	Broker  *Broker
	Session models.Session
}

func (wsc *WSConn) Send(from string, to string, topic string, body any) {
	data, err := json.Marshal(map[string]any{
		"type": topic,
		"body": body,
		"from": from,
	})
	if err != nil {
		log.Printf("Error marshal: %v", err)
		return
	}

	wsc.Msg <- data
}

func (wsc *WSConn) ReadLoop() {
	defer func() {
		wsc.Broker.RemoveService(wsc.ID)
		wsc.Conn.Close()
	}()

	wsc.Conn.SetReadLimit(1024)
	wsc.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	wsc.Conn.SetPongHandler(func(string) error {
		wsc.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		return nil
	})

	for {
		var msg IncomingWebSocketMessage
		err := wsc.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Handle Message ...
		log.Printf("Message: %+v", msg)

		switch msg.Type {
		case "newGame":
			switch {
			case msg.Voidsent != nil:
				log.Printf("Start a new match of Voidsent: %+v", *msg.Voidsent)
			}
		}
	}
}

func (wsc *WSConn) WriteLoop() {
	ticker := time.NewTicker(20 * time.Second)
	defer func() {
		ticker.Stop()
		wsc.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-wsc.Msg:
			log.Printf("Sending msg: %s", msg)
			wsc.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Channel is closed
				return
			}

			log.Printf("Write Message: %v", msg)
			wsc.Conn.WriteMessage(websocket.TextMessage, msg)

		case <-ticker.C:
			wsc.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := wsc.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping Error: %v", err)
				return
			}
		}
	}
}

// func (wsc *WSConn) newGame(s *melody.Sessions, session *models.Sessions, newGameMsg *newGameMessage) {
// 	game := &voidsent.Game{
// 		ID:       ulid.MustNew(ulid.Now(), rand.Reader).String(),
// 		Status:   voidsent.StatusAwaitingPlayer,
// 		Name:     newGameMsg.Name,
// 		Melody:   ws.Melody,
// 		Sessons: make([]string, 0),
// 	}
//
// 	ws.Hub.CreateRoom(session.ID, game)
// }

func (wsc *WSConn) Disconnect() {
	log.Println("Disconnected")
}

func (wsc *WSConn) chat(topic string, msg string) {
	switch topic {
	case "room:join":
		data, err := json.Marshal(map[string]any{
			"type": "join",
			"room": msg,
			"name": wsc.Session.Username,
		})
		if err != nil {
			log.Println("Error marshal", err)
		}
		wsc.Msg <- data
	}
}
