package services

type Service interface {
	Send(from string, to string, topic string, body any)
}
