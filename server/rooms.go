package server

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	ROOM_WAITING  = "waiting"
	ROOM_STARTING = "starting"
	ROOM_PLAYING  = "playing"
	ROOM_FINISHED = "finished" // Nuevo estado
)

type Room struct {
	ID         string             `json:"id"`
	Players    map[string]*Player `json:"players"`
	GameState  GameState          `json:"game_state"`
	State      string             `json:"state"`
	Mutex      sync.Mutex         `json:"-"`
	lastUpdate time.Time          `json:"-"`
	StartTime  time.Time          `json:"-"`
	Winner     string             `json:"winner"`    // ID del jugador ganador
	WinReason  string             `json:"winReason"` // Razón de la victoria
}

func NewRoom(id string) *Room {
	room := &Room{
		ID:         id,
		Players:    make(map[string]*Player),
		State:      ROOM_WAITING,
		StartTime:  time.Now(),
		lastUpdate: time.Now(),
		GameState: GameState{
			Paddle1: Paddle{X: PADDLE1_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
			Paddle2: Paddle{X: PADDLE2_X, Y: FIELD_HEIGHT/2 - PADDLE_HEIGHT/2, Width: PADDLE_WIDTH, Height: PADDLE_HEIGHT},
			Balls:   []Ball{}, // Sin pelotas hasta que empiece el juego
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
		r.State = ROOM_WAITING
	} else {
		p.Role = "player2"
		r.State = ROOM_STARTING
		r.StartTime = time.Now().Add(3 * time.Second)
		log.Printf("Second player joined. Starting game in 3 seconds. Current state: %s", r.State)
	}

	r.Players[p.ID] = p
	log.Printf("Player %s joined room %s as %s. Room state: %s", p.ID, r.ID, p.Role, r.State)
	return nil
}

func (r *Room) PlayerDisconnected(playerID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	disconnectedPlayer := r.Players[playerID]
	if r.State == ROOM_PLAYING {
		// Si el jugador 1 se desconecta, gana el jugador 2 y viceversa
		if disconnectedPlayer.Role == "player1" {
			r.Winner = "player2"
		} else {
			r.Winner = "player1"
		}
		r.WinReason = "opponent_disconnected"
		r.State = ROOM_FINISHED
		log.Printf("Game finished. Player %s wins by disconnection", r.Winner)
	}
	delete(r.Players, playerID)
}
