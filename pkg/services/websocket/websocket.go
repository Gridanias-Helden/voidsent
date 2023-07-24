package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gridanias-helden/voidsent/pkg/middleware"
	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/services"
	"github.com/gridanias-helden/voidsent/pkg/storage"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type WebSocket struct {
	Sessons storage.Sessions
	Broker  *services.Broker
}

func (ws *WebSocket) HTTPRequest(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionKey).(models.Session)
	if !ok {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	wsConn := &services.WSConn{
		ID:      "ws:" + sess.ID,
		Broker:  ws.Broker,
		Conn:    conn,
		Msg:     make(chan []byte),
		Session: sess,
	}
	go wsConn.ReadLoop()
	go wsConn.WriteLoop()
	time.Sleep(50 * time.Millisecond)

	ws.Broker.AddService(wsConn.ID, wsConn)
	ws.Broker.Send(wsConn.ID, "chat", "lobby:join", wsConn)
}
