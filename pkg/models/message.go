package models

type Message struct {
	To    string `json:"to"`
	From  string `json:"from"`
	Topic string `json:"topic"`
	Body  any    `json:"body"`
}
