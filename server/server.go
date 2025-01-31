package server

import (
    "sync"
)

type Server struct {
    Rooms map[string]*Room
    Mutex sync.Mutex
}

func NewServer() *Server {
    return &Server{
        Rooms: make(map[string]*Room),
    }
}
