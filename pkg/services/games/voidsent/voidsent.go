package voidsent

import (
	"fmt"

	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/services"
)

type voidsent struct {
	in chan models.Message
	m  *services.Broker
}

func New(m *services.Broker) services.Service {
	v := &voidsent{
		in: make(chan models.Message),
		m:  m,
	}

	v.Start()

	return v
}

func (v *voidsent) Send(from string, to string, topic string, body any) {
	v.in <- models.Message{To: to, From: from, Body: body}
}

func (v *voidsent) Start() {
	go func() {
		for msg := range v.in {
			fmt.Printf("Log: %+v\n", msg)
		}
	}()
}
