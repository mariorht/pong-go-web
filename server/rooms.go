package server

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Room struct {
	ID        string             `json:"id"`
	Players   map[string]*Player `json:"players"`
	GameState GameState          `json:"game_state"`
	Mutex     sync.Mutex         `json:"-"`
	engine    PhysicsEngine      `json:"-"`
	StartTime time.Time          `json:"-"`
}

func NewRoom(id string) *Room {
	balls := make([]Ball, 2)
	for i := range balls {
		balls[i] = createNewBall()
	}

	room := &Room{
		ID:      id,
		Players: make(map[string]*Player),
		GameState: GameState{
			Paddle1: Paddle{X: PADDLE1_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
			Paddle2: Paddle{X: PADDLE2_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
			Balls:   balls,
		},
	}

	// Inicializar el engine después de crear la room
	room.engine = PhysicsEngine{
		room:       room,
		lastUpdate: time.Now(),
	}

	return room
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
