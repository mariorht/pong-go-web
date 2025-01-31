package server

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/google/uuid"
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

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

// Manejo de conexiones WebSocket con logs
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading connection:", err)
        http.Error(w, "Could not open WebSocket connection", http.StatusInternalServerError)
        return
    }

    log.Printf("New client connected from %s", r.RemoteAddr) // Log de nueva conexión

    clientID := uuid.New().String()
    player := &Player{Conn: conn, ID: clientID}
    roomID := "default" // Assign to a default room for simplicity
    s.Mutex.Lock()
    if _, ok := s.Rooms[roomID]; !ok {
        s.Rooms[roomID] = NewRoom(roomID)
    }
    s.Rooms[roomID].AddPlayer(player)
    log.Printf("Player added to room %s: %v", roomID, player)
    s.Mutex.Unlock()

    err = conn.WriteJSON(map[string]string{"type": "client_id", "id": clientID})
    if err != nil {
        log.Printf("Error sending client ID to %s: %v", r.RemoteAddr, err)
        return
    }

    defer func() {
        log.Printf("Client disconnected from %s", r.RemoteAddr) // Log de desconexión
        s.Mutex.Lock()
        delete(s.Rooms[roomID].Players, player.ID)
        log.Printf("Player removed from room %s: %v", roomID, player)
        s.Mutex.Unlock()
        conn.Close()
    }()

    // Escuchar mensajes del cliente
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message from %s: %v", r.RemoteAddr, err)
            break
        }
        log.Printf("Received message from %s: %s", r.RemoteAddr, string(msg))
    }
}

func (s *Server) GetServerStatus(w http.ResponseWriter, r *http.Request) {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()



    status := map[string]interface{}{
        "rooms":   s.Rooms,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(status); err != nil {
        log.Println("Error encoding server status:", err)
        http.Error(w, "Could not encode server status", http.StatusInternalServerError)
    }
}
