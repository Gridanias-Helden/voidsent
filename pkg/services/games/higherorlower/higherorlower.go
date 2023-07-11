package higherorlower

import (
	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/services"
)

type higherOrLower struct {
	in chan models.Message
	m  *services.Broker

	max   int
	min   int
	tries int
	value int
}

func New(m *services.Broker) services.Service {
	hol := &higherOrLower{
		in:    make(chan models.Message),
		m:     m,
		max:   100,
		min:   1,
		tries: 0,
	}

	hol.Start()

	return hol
}

func (hol *higherOrLower) Send(from string, to string, topic string, body any) {
	hol.in <- models.Message{
		To:    to,
		From:  from,
		Topic: topic,
		Body:  body,
	}
}

func (hol *higherOrLower) Start() {
	go func() {
		for msg := range hol.in {
			t, ok := msg.Body.(string)
			if !ok {
				continue
			}

			switch t {
			case "start":
				hol.min = 1
				hol.max = 100
				hol.tries = 1
				hol.value = (hol.min + hol.max) / 2
				hol.m.Send(msg.To, msg.From, "", map[string]int{"tries": hol.tries, "value": hol.value})

			case "higher":
				hol.min = hol.value + 1
				hol.value = (hol.min + hol.max) / 2
				hol.tries++
				hol.m.Send(msg.To, msg.From, "", map[string]int{"tries": hol.tries, "value": hol.value})

			case "lower":
				hol.max = hol.value - 1
				hol.value = (hol.min + hol.max) / 2
				hol.tries++
				hol.m.Send(msg.To, msg.From, "", map[string]int{"tries": hol.tries, "value": hol.value})

			case "exit":
				hol.m.Send(msg.To, msg.From, "exit", "goodbye")
				return
			}
		}
	}()
}
