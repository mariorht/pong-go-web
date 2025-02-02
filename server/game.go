package server

import (
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	// Dimensiones del campo
	FIELD_WIDTH  = 1000
	FIELD_HEIGHT = 600

	// Dimensiones de las palas
	PADDLE_WIDTH  = 10
	PADDLE_HEIGHT = 100
	PADDLE1_X     = 50
	PADDLE2_X     = FIELD_WIDTH - 50 - PADDLE_WIDTH

	// Dimensiones de la pelota
	BALL_RADIUS  = 10
	BALL_START_X = FIELD_WIDTH / 2
	BALL_START_Y = FIELD_HEIGHT / 2

	// Velocidades
	BASE_BALL_SPEED      = 5.0
	BALL_SPEED_VARIATION = 2.0
	PADDLE_SPEED         = 10
)

type GameState struct {
	Paddle1  Paddle `json:"paddle1"`
	Paddle2  Paddle `json:"paddle2"`
	Balls    []Ball `json:"balls"`
	Score1   int    `json:"score1"`
	Score2   int    `json:"score2"`
	GameTime int    `json:"gameTime"` // Tiempo en segundos
}

type Paddle struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Ball struct {
	X      float64 `json:"-"` // Uso interno
	Y      float64 `json:"-"` // Uso interno
	Radius int     `json:"radius"`
	VX     float64 `json:"-"` // Uso interno
	VY     float64 `json:"-"` // Uso interno
}

// Nueva estructura para la vista
type BallView struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Radius int `json:"radius"`
	VX     int `json:"vx"`
	VY     int `json:"vy"`
}

func (b Ball) ToView() BallView {
	return BallView{
		X:      int(b.X),
		Y:      int(b.Y),
		Radius: b.Radius,
		VX:     int(b.VX),
		VY:     int(b.VY),
	}
}

type PhysicsEngine struct {
	room *Room
}

func (s *Server) StartGameLoop() {
	physicsUpdate := time.NewTicker(10 * time.Millisecond) // 100 FPS para física
	renderUpdate := time.NewTicker(20 * time.Millisecond)  // 50 FPS para renderizado
	gameStartTime := time.Now()
	lastBallTime := time.Now()
	defer physicsUpdate.Stop()
	defer renderUpdate.Stop()

	// Iniciar el hilo de física
	go func() {
		for range physicsUpdate.C {
			s.Mutex.Lock()
			for _, room := range s.Rooms {
				engine := PhysicsEngine{room: room}
				engine.updatePhysics()
			}
			s.Mutex.Unlock()
		}
	}()

	// Hilo principal para renderizado
	for range renderUpdate.C {
		s.Mutex.Lock()
		for _, room := range s.Rooms {
			currentTime := time.Now()
			room.GameState.GameTime = int(currentTime.Sub(gameStartTime).Seconds())

			// Añadir nueva pelota cada 10 segundos
			if currentTime.Sub(lastBallTime).Seconds() >= 10 {
				room.GameState.Balls = append(room.GameState.Balls, Ball{
					X:      BALL_START_X,
					Y:      BALL_START_Y,
					Radius: BALL_RADIUS,
					VX:     rand.Float64()*10 - 5,
					VY:     rand.Float64()*10 - 5,
				})
				lastBallTime = currentTime
			}

			s.broadcastGameState(room)
		}
		s.Mutex.Unlock()
	}
}

func (engine *PhysicsEngine) updatePhysics() {
	room := engine.room

	for i := range room.GameState.Balls {
		ball := &room.GameState.Balls[i]

		// Update ball position with floating point precision
		ball.X += ball.VX
		ball.Y += ball.VY

		// Check for collisions with top and bottom walls
		if ball.Y-float64(ball.Radius) <= 0 || ball.Y+float64(ball.Radius) >= FIELD_HEIGHT {
			ball.VY = -ball.VY
		}

		// Check for collisions with paddles
		// Paddle 1
		if ball.X-float64(ball.Radius) <= float64(room.GameState.Paddle1.X+PADDLE_WIDTH) &&
			ball.Y >= float64(room.GameState.Paddle1.Y) &&
			ball.Y <= float64(room.GameState.Paddle1.Y+PADDLE_HEIGHT) {
			ball.VX = -ball.VX
		}
		// Paddle 2
		if ball.X+float64(ball.Radius) >= float64(room.GameState.Paddle2.X) &&
			ball.Y >= float64(room.GameState.Paddle2.Y) &&
			ball.Y <= float64(room.GameState.Paddle2.Y+PADDLE_HEIGHT) {
			ball.VX = -ball.VX
		}

		// Check for goals
		if ball.X-float64(ball.Radius) <= 0 {
			room.GameState.Score2++
			resetBall(ball)
		} else if ball.X+float64(ball.Radius) >= FIELD_WIDTH {
			room.GameState.Score1++
			resetBall(ball)
		}
	}
}

func resetBall(ball *Ball) {
	ball.X = BALL_START_X
	ball.Y = BALL_START_Y
	// Velocidades aleatorias más precisas
	angle := rand.Float64() * 2 * math.Pi
	speed := BASE_BALL_SPEED + rand.Float64()*BALL_SPEED_VARIATION // Velocidad base 5-7
	ball.VX = speed * math.Cos(angle)
	ball.VY = speed * math.Sin(angle)
}

// Añadir esta estructura cerca de las otras definiciones de tipos
type GameConfig struct {
	FieldWidth   int `json:"fieldWidth"`
	FieldHeight  int `json:"fieldHeight"`
	PaddleWidth  int `json:"paddleWidth"`
	PaddleHeight int `json:"paddleHeight"`
	BallRadius   int `json:"ballRadius"`
}

func (s *Server) broadcastGameState(room *Room) {
	// Convertir las pelotas a su versión de vista
	ballViews := make([]BallView, len(room.GameState.Balls))
	for i, ball := range room.GameState.Balls {
		ballViews[i] = ball.ToView()
	}

	viewState := struct {
		Paddle1  Paddle     `json:"paddle1"`
		Paddle2  Paddle     `json:"paddle2"`
		Balls    []BallView `json:"balls"`
		Score1   int        `json:"score1"`
		Score2   int        `json:"score2"`
		GameTime int        `json:"gameTime"`
	}{
		Paddle1:  room.GameState.Paddle1,
		Paddle2:  room.GameState.Paddle2,
		Balls:    ballViews,
		Score1:   room.GameState.Score1,
		Score2:   room.GameState.Score2,
		GameTime: room.GameState.GameTime,
	}

	sendTime := time.Now().UnixNano() / int64(time.Microsecond)
	message := map[string]interface{}{
		"type":     "game_state",
		"state":    viewState,
		"sendTime": sendTime,
		"config": GameConfig{
			FieldWidth:   FIELD_WIDTH,
			FieldHeight:  FIELD_HEIGHT,
			PaddleWidth:  PADDLE_WIDTH,
			PaddleHeight: PADDLE_HEIGHT,
			BallRadius:   BALL_RADIUS,
		},
	}

	for _, player := range room.Players {
		player.mutex.Lock()
		if err := player.Conn.WriteJSON(message); err != nil {
			log.Printf("Error sending game state to %s: %v", player.ID, err)
		}
		player.mutex.Unlock()
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
	const minY = 0
	const maxY = FIELD_HEIGHT - PADDLE_HEIGHT

	if player.Role == "player1" {
		paddle := &room.GameState.Paddle1
		if direction == "ArrowUp" {
			newY := paddle.Y - PADDLE_SPEED
			paddle.Y = int(math.Max(float64(minY), float64(newY)))
		} else if direction == "ArrowDown" {
			newY := paddle.Y + PADDLE_SPEED
			paddle.Y = int(math.Min(float64(maxY), float64(newY)))
		}
	} else if player.Role == "player2" {
		paddle := &room.GameState.Paddle2
		if direction == "ArrowUp" {
			newY := paddle.Y - PADDLE_SPEED
			paddle.Y = int(math.Max(float64(minY), float64(newY)))
		} else if direction == "ArrowDown" {
			newY := paddle.Y + PADDLE_SPEED
			paddle.Y = int(math.Min(float64(maxY), float64(newY)))
		}
	}
}
