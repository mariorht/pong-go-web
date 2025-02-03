package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn        *websocket.Conn `json:"-"`
	ID          string          `json:"id"`
	Role        string          `json:"role"`
	Name        string          `json:"name"` // Nuevo campo para el nombre del jugador
	mutex       sync.Mutex      `json:"-"`
	isConnected bool
}
