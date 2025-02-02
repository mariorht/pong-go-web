package server

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
)

type Room struct {
	ID        string             `json:"id"`
	Players   map[string]*Player `json:"players"`
	GameState GameState          `json:"game_state"`
	Mutex     sync.Mutex         `json:"-"`
}

func NewRoom(id string) *Room {
	balls := make([]Ball, 100)
	for i := range balls {
		angle := rand.Float64() * 2 * math.Pi
		speed := 5.0 + rand.Float64()*2 // Velocidad base 5-7
		balls[i] = Ball{
			X:      400,
			Y:      200,
			Radius: 10,
			VX:     speed * math.Cos(angle),
			VY:     speed * math.Sin(angle),
		}
	}

	return &Room{
		ID:      id,
		Players: make(map[string]*Player),
		GameState: GameState{
			Paddle1: Paddle{X: 50, Y: 150, Width: 10, Height: 100},
			Paddle2: Paddle{X: 740, Y: 150, Width: 10, Height: 100},
			Balls:   balls,
		},
	}
}

func (r *Room) AddPlayer(p *Player) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if len(r.Players) >= 2 {
		return fmt.Errorf("room %s is full", r.ID)
	}

	if len(r.Players) == 0 {
		p.Role = "player1"
	} else {
		p.Role = "player2"
	}

	r.Players[p.ID] = p
	log.Printf("Player %s joined room %s as %s", p.ID, r.ID, p.Role)
	return nil
}
