package server

import "github.com/gorilla/websocket"

type Player struct {
    Conn *websocket.Conn `json:"-"`
    ID   string          `json:"id"`
    Role string          `json:"role"`
}
