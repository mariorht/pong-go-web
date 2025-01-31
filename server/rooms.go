package server

import (
    "log"
    "sync"
)

type Room struct {
    ID        string            `json:"id"`
    Players   map[string]*Player `json:"players"`
    GameState GameState         `json:"game_state"`
    Mutex     sync.Mutex        `json:"-"`
}

func NewRoom(id string) *Room {
    return &Room{
        ID: id,
        Players: make(map[string]*Player),
        GameState: GameState{
            Paddle1: Paddle{X: 50, Y: 150, Width: 10, Height: 100},
            Paddle2: Paddle{X: 740, Y: 150, Width: 10, Height: 100},
            Ball:    Ball{X: 400, Y: 200, Radius: 10, VX: 5, VY: 5},
        },
    }
}

func (r *Room) AddPlayer(p *Player) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()
    
    r.Players[p.ID] = p
    log.Printf("Player joined room %s", r.ID)
}
