package services

//import (
//	"log"
//)

type Broker struct {
	services map[string]Service
}

func NewBroker() *Broker {
	return &Broker{
		services: make(map[string]Service),
	}
}

func (m *Broker) Send(from string, to string, topic string, body any) {
	//	log.Printf("From: %q, To: %q, Topic: %q, Content: %v", from, to, topic, body)

	if actor, ok := m.services[to]; ok {
		go actor.Send(from, to, topic, body)
	}
}

func (m *Broker) AddService(name string, actor Service) {
	m.services[name] = actor
}

func (m *Broker) RemoveService(name string) {
	delete(m.services, name)
}
