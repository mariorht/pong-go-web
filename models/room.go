package models

import "sync"

type Room struct {
    ID        string            `json:"id"`
    Players   map[string]*Player `json:"players"`
    GameState GameState         `json:"game_state"`
    Mutex     sync.Mutex        `json:"-"`
}
