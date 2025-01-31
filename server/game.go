package server

import (
    "time"
    "log"
)

type GameState struct {
    Paddle1 Paddle `json:"paddle1"`
    Paddle2 Paddle `json:"paddle2"`
    Ball    Ball   `json:"ball"`
}

type Paddle struct {
    X      int `json:"x"`
    Y      int `json:"y"`
    Width  int `json:"width"`
    Height int `json:"height"`
}

type Ball struct {
    X      int `json:"x"`
    Y      int `json:"y"`
    Radius int `json:"radius"`
    VX     int `json:"vx"`
    VY     int `json:"vy"`
}

func (s *Server) StartGameLoop() {
    ticker := time.NewTicker(50 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        s.Mutex.Lock()
        for _, room := range s.Rooms {
            s.updateGameState(room)
            s.broadcastGameState(room)
        }
        s.Mutex.Unlock()
    }
}

func (s *Server) updateGameState(room *Room) {
    // Update ball position
    room.GameState.Ball.X += room.GameState.Ball.VX
    room.GameState.Ball.Y += room.GameState.Ball.VY

    // Check for collisions with top and bottom walls
    if room.GameState.Ball.Y-room.GameState.Ball.Radius <= 0 || room.GameState.Ball.Y+room.GameState.Ball.Radius >= 400 {
        room.GameState.Ball.VY = -room.GameState.Ball.VY
    }

    // Check for collisions with paddles
    // Paddle 1
    if room.GameState.Ball.X-room.GameState.Ball.Radius <= room.GameState.Paddle1.X+room.GameState.Paddle1.Width &&
        room.GameState.Ball.Y >= room.GameState.Paddle1.Y &&
        room.GameState.Ball.Y <= room.GameState.Paddle1.Y+room.GameState.Paddle1.Height {
        room.GameState.Ball.VX = -room.GameState.Ball.VX
    }
    // Paddle 2
    if room.GameState.Ball.X+room.GameState.Ball.Radius >= room.GameState.Paddle2.X &&
        room.GameState.Ball.Y >= room.GameState.Paddle2.Y &&
        room.GameState.Ball.Y <= room.GameState.Paddle2.Y+room.GameState.Paddle2.Height {
        room.GameState.Ball.VX = -room.GameState.Ball.VX
    }
}

func (s *Server) broadcastGameState(room *Room) {
    gameState := room.GameState

    for _, player := range room.Players {
        if err := player.Conn.WriteJSON(map[string]interface{}{"type": "game_state", "state": gameState}); err != nil {
            log.Printf("Error sending game state to %s: %v", player.ID, err)
        }
    }
}

func (s *Server) handlePlayerMove(roomID, playerID, direction string) {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()

    room, ok := s.Rooms[roomID]
    if !ok {
        return
    }

    // Update paddle position based on player input
    if playerID == "player1" {
        if direction == "ArrowUp" {
            room.GameState.Paddle1.Y -= 10
        } else if direction == "ArrowDown" {
            room.GameState.Paddle1.Y += 10
        }
    } else if playerID == "player2" {
        if direction == "ArrowUp" {
            room.GameState.Paddle2.Y -= 10
        } else if direction == "ArrowDown" {
            room.GameState.Paddle2.Y += 10
        }
    }
}
