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
	BALL_RADIUS  = 20
	BALL_START_X = FIELD_WIDTH / 2
	BALL_START_Y = FIELD_HEIGHT / 2

	// Velocidades
	BASE_BALL_SPEED      = 0.5
	BALL_SPEED_VARIATION = 0.1
	PADDLE_SPEED         = 20
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
	physicsUpdate := time.NewTicker(2 * time.Millisecond) // 500 FPS para física
	renderUpdate := time.NewTicker(20 * time.Millisecond) // 50 FPS para renderizado
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
				room.GameState.Balls = append(room.GameState.Balls, createNewBall())
				lastBallTime = currentTime
			}

			s.broadcastGameState(room)
		}
		s.Mutex.Unlock()
	}
}

func (engine *PhysicsEngine) updatePhysics() {
	room := engine.room

	// Primero actualizar posiciones
	for i := len(room.GameState.Balls) - 1; i >= 0; i-- {
		ball := &room.GameState.Balls[i]
		ball.X += ball.VX
		ball.Y += ball.VY

		// Colisiones con paredes y palas (código existente)
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

		// Check for goals and remove ball
		if ball.X-float64(ball.Radius) <= 0 {
			room.GameState.Score2++
			// Eliminar la pelota
			room.GameState.Balls = append(room.GameState.Balls[:i], room.GameState.Balls[i+1:]...)
		} else if ball.X+float64(ball.Radius) >= FIELD_WIDTH {
			room.GameState.Score1++
			// Eliminar la pelota
			room.GameState.Balls = append(room.GameState.Balls[:i], room.GameState.Balls[i+1:]...)
		}
	}

	// Después comprobar colisiones entre pelotas
	for i := 0; i < len(room.GameState.Balls); i++ {
		for j := i + 1; j < len(room.GameState.Balls); j++ {
			ball1 := &room.GameState.Balls[i]
			ball2 := &room.GameState.Balls[j]

			if checkBallCollision(ball1, ball2) {
				resolveBallCollision(ball1, ball2)
			}
		}
	}

	// Si no quedan pelotas, crear una nueva
	if len(room.GameState.Balls) == 0 {
		room.GameState.Balls = append(room.GameState.Balls, createNewBall())
	}
}

// Detectar si dos pelotas están colisionando
func checkBallCollision(ball1, ball2 *Ball) bool {
	dx := ball1.X - ball2.X
	dy := ball1.Y - ball2.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < float64(ball1.Radius+ball2.Radius)
}

// Resolver la colisión elástica entre dos pelotas
func resolveBallCollision(ball1, ball2 *Ball) {
	// Vector normal de colisión
	nx := ball2.X - ball1.X
	ny := ball2.Y - ball1.Y

	// Distancia entre centros
	d := math.Sqrt(nx*nx + ny*ny)
	if d == 0 {
		return // Evitar división por cero
	}

	// Normalizar el vector normal
	nx = nx / d
	ny = ny / d

	// Velocidad relativa
	dvx := ball1.VX - ball2.VX
	dvy := ball1.VY - ball2.VY

	// Velocidad relativa en dirección normal
	velAlongNormal := dvx*nx + dvy*ny

	// No colisionar si las pelotas se están alejando
	if velAlongNormal > 0 {
		return
	}

	// Coeficiente de restitución (1 para colisión elástica perfecta)
	restitution := 1.0

	// Impulso
	j := -(1 + restitution) * velAlongNormal
	j = j / 2 // Dividir por 2 porque ambas pelotas tienen la misma masa

	// Aplicar impulso
	ball1.VX += j * nx
	ball1.VY += j * ny
	ball2.VX -= j * nx
	ball2.VY -= j * ny

	// Separar las pelotas para evitar que se peguen
	overlap := float64(ball1.Radius+ball2.Radius) - d
	if overlap > 0 {
		separation := overlap / 2
		ball1.X -= nx * separation
		ball1.Y -= ny * separation
		ball2.X += nx * separation
		ball2.Y += ny * separation
	}
}

func createNewBall() Ball {
	angle := rand.Float64() * 2 * math.Pi
	speed := BASE_BALL_SPEED + rand.Float64()*BALL_SPEED_VARIATION
	return Ball{
		X:      BALL_START_X,
		Y:      BALL_START_Y,
		Radius: BALL_RADIUS,
		VX:     speed * math.Cos(angle),
		VY:     speed * math.Sin(angle),
	}
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

	// Iterar de forma segura sobre los jugadores y limpiar los desconectados
	disconnectedPlayers := []string{}

	for id, player := range room.Players {
		if !player.isConnected {
			disconnectedPlayers = append(disconnectedPlayers, id)
			continue
		}

		player.mutex.Lock()
		if err := player.Conn.WriteJSON(message); err != nil {
			log.Printf("Error sending game state to %s: %v", player.ID, err)
			player.isConnected = false
			disconnectedPlayers = append(disconnectedPlayers, id)
		}
		player.mutex.Unlock()
	}

	// Limpiar jugadores desconectados
	for _, id := range disconnectedPlayers {
		delete(room.Players, id)
		log.Printf("Removed disconnected player: %s", id)
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
