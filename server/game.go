package server

import (
	"log"
	"time"
)

type GameState struct {
	Paddle1 Paddle `json:"paddle1"`
	Paddle2 Paddle `json:"paddle2"`
	Balls   []Ball `json:"balls"`
	Score1  int    `json:"score1"`
	Score2  int    `json:"score2"`
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

//Versión concurrente, es bastante más lenta que la versión secuencial
//
// func (s *Server) updateGameState(room *Room) {
// 	start := time.Now()

// 	var wg sync.WaitGroup

// 	for i := range room.GameState.Balls {
// 		wg.Add(1)
// 		go func(ball *Ball) {
// 			defer wg.Done()

// 			// Update ball position
// 			ball.X += ball.VX
// 			ball.Y += ball.VY

// 			// Check for collisions with top and bottom walls
// 			if ball.Y-ball.Radius <= 0 || ball.Y+ball.Radius >= 400 {
// 				ball.VY = -ball.VY
// 			}

// 			// Check for collisions with paddles
// 			// Paddle 1
// 			if ball.X-ball.Radius <= room.GameState.Paddle1.X+room.GameState.Paddle1.Width &&
// 				ball.Y >= room.GameState.Paddle1.Y &&
// 				ball.Y <= room.GameState.Paddle1.Y+room.GameState.Paddle1.Height {
// 				ball.VX = -ball.VX
// 			}
// 			// Paddle 2
// 			if ball.X+ball.Radius >= room.GameState.Paddle2.X &&
// 				ball.Y >= room.GameState.Paddle2.Y &&
// 				ball.Y <= room.GameState.Paddle2.Y+room.GameState.Paddle2.Height {
// 				ball.VX = -ball.VX
// 			}

// 			// Check for goals
// 			if ball.X-ball.Radius <= 0 {
// 				room.GameState.Score2++
// 				s.resetBall(ball)
// 			} else if ball.X+ball.Radius >= 800 {
// 				room.GameState.Score1++
// 				s.resetBall(ball)
// 			}
// 		}(&room.GameState.Balls[i])
// 	}

// 	wg.Wait()

// 	elapsed := time.Since(start)
// 	log.Printf("updateGameState took %s", elapsed)
// }

func (s *Server) updateGameState(room *Room) {
	// start := time.Now()

	for i := range room.GameState.Balls {
		ball := &room.GameState.Balls[i]

		// Update ball position
		ball.X += ball.VX
		ball.Y += ball.VY

		// Check for collisions with top and bottom walls
		if ball.Y-ball.Radius <= 0 || ball.Y+ball.Radius >= 400 {
			ball.VY = -ball.VY
		}

		// Check for collisions with paddles
		// Paddle 1
		if ball.X-ball.Radius <= room.GameState.Paddle1.X+room.GameState.Paddle1.Width &&
			ball.Y >= room.GameState.Paddle1.Y &&
			ball.Y <= room.GameState.Paddle1.Y+room.GameState.Paddle1.Height {
			ball.VX = -ball.VX
		}
		// Paddle 2
		if ball.X+ball.Radius >= room.GameState.Paddle2.X &&
			ball.Y >= room.GameState.Paddle2.Y &&
			ball.Y <= room.GameState.Paddle2.Y+room.GameState.Paddle2.Height {
			ball.VX = -ball.VX
		}

		// Check for goals
		if ball.X-ball.Radius <= 0 {
			room.GameState.Score2++
			s.resetBall(ball)
		} else if ball.X+ball.Radius >= 800 {
			room.GameState.Score1++
			s.resetBall(ball)
		}
	}

	// elapsed := time.Since(start)
	// log.Printf("updateGameState took %s", elapsed)
}

func (s *Server) resetBall(ball *Ball) {
	ball.X = 400
	ball.Y = 200
	ball.VX = 5
	ball.VY = 5
}

func (s *Server) broadcastGameState(room *Room) {
	gameState := room.GameState
	sendTime := time.Now().UnixNano() / int64(time.Microsecond)

	for _, player := range room.Players {
		if err := player.Conn.WriteJSON(map[string]interface{}{"type": "game_state", "state": gameState, "sendTime": sendTime}); err != nil {
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

	player := room.Players[playerID]

	// Update paddle position based on player input
	if player.Role == "player1" {
		if direction == "ArrowUp" {
			room.GameState.Paddle1.Y -= 10
		} else if direction == "ArrowDown" {
			room.GameState.Paddle1.Y += 10
		}
	} else if player.Role == "player2" {
		if direction == "ArrowUp" {
			room.GameState.Paddle2.Y -= 10
		} else if direction == "ArrowDown" {
			room.GameState.Paddle2.Y += 10
		}
	}
}
