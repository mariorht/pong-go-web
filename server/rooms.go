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
	balls := make([]Ball, 1)
	for i := range balls {
		angle := rand.Float64() * 2 * math.Pi
		speed := BASE_BALL_SPEED + rand.Float64()*BALL_SPEED_VARIATION
		balls[i] = Ball{
			X:      BALL_START_X,
			Y:      BALL_START_Y,
			Radius: BALL_RADIUS,
			VX:     speed * math.Cos(angle),
			VY:     speed * math.Sin(angle),
		}
	}

	return &Room{
		ID:      id,
		Players: make(map[string]*Player),
		GameState: GameState{
			Paddle1: Paddle{X: PADDLE1_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
			Paddle2: Paddle{X: PADDLE2_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
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
