package server

import (
    "log"
    "sync"
)

type Room struct {
    ID      string            `json:"id"`
    Players map[string]*Player `json:"players"`
    Mutex   sync.Mutex        `json:"-"`
}

func NewRoom(id string) *Room {
    return &Room{
        ID:      id,
        Players: make(map[string]*Player),
    }
}

// Agrega un jugador a la sala
func (r *Room) AddPlayer(p *Player) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    
    r.Players[p.ID] = p
    log.Printf("Player joined room %s", r.ID)
}
