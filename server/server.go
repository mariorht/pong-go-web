package server

import (
	"sync"
	"time"
)

type Server struct {
	Rooms         map[string]*Room
	Mutex         sync.Mutex
	gameStartTime time.Time
	lastBallTime  time.Time
}

func NewServer() *Server {
	return &Server{
		Rooms:         make(map[string]*Room),
		gameStartTime: time.Now(),
		lastBallTime:  time.Now(),
	}
}
