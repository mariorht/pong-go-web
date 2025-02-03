package server

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Room struct {
	ID         string             `json:"id"`
	Players    map[string]*Player `json:"players"`
	GameState  GameState          `json:"game_state"`
	Mutex      sync.Mutex         `json:"-"`
	lastUpdate time.Time          `json:"-"`
	StartTime  time.Time          `json:"-"`
}

func NewRoom(id string) *Room {
	room := &Room{
		ID:         id,
		Players:    make(map[string]*Player),
		StartTime:  time.Now(),
		lastUpdate: time.Now(),
		GameState: GameState{
			Paddle1: Paddle{X: PADDLE1_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
			Paddle2: Paddle{X: PADDLE2_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
			Balls:   []Ball{createNewBall()},
		},
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
