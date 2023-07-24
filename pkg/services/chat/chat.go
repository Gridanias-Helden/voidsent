package chat

import (
	"log"

	"github.com/gridanias-helden/voidsent/pkg/services"
)

type chatHandler struct {
	users  map[string]*services.WSConn
	rooms  map[string]string
	in     chan Message
	broker *services.Broker
}

func New(broker *services.Broker) services.Service {
	return &chatHandler{
		users:  make(map[string]*services.WSConn, 0),
		rooms:  make(map[string]string),
		broker: broker,
		in:     make(chan Message),
	}
}

func (c *chatHandler) Send(from string, to string, topic string, body any) {
	log.Printf("Chat: From %q, To: %q, Topic: %q, Body: %+v", from, to, topic, body)

	switch topic {
	case "lobby:join":
		sess, ok := body.(*services.WSConn)
		if !ok {
			log.Printf("No valid connection %t", body)
			return
		}

		c.users[from] = sess
		c.rooms[from] = "lobby"

		c.publish("lobby", to, "room:join", map[string]any{"room": "lobby", "name": sess.Session.Username, "avatar": sess.Session.Avatar})
		c.broker.Send(to, from, "session", map[string]any{"name": sess.Session.Username, "avatar": sess.Session.Avatar})

	case "room:join":
		room, ok := body.(string)
		if !ok {
			log.Printf("room should be a string")
			return
		}

		if _, ok := c.rooms[from]; !ok {
			log.Printf("%q is not any room", from)
			return
		}

		if c.rooms[from] != "lobby" {
			log.Printf("Can only join a room from the lobby")
			return
		}

		c.rooms[from] = room
		c.broker.Send(to, from, topic, room)
	}
}

func (c *chatHandler) publish(room string, from string, topic string, msg any) {
	for to := range c.users {
		if c.rooms[to] == room {
			c.broker.Send(from, to, topic, msg)
		}
	}
}
