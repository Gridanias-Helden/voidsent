package models

type LobbyEntry struct {
	ID      string `json:"id"`
	Game    string `json:"game"`
	Name    string `json:"name"`
	Players int    `json:"players"`
}
