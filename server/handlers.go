package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusInternalServerError)
		return
	}

	log.Printf("New client connected from %s", r.RemoteAddr)

	clientID := uuid.New().String()
	player := &Player{Conn: conn, ID: clientID}
	roomID := "default"
	s.Mutex.Lock()
	if _, ok := s.Rooms[roomID]; !ok {
		s.Rooms[roomID] = NewRoom(roomID)
	}
	err = s.Rooms[roomID].AddPlayer(player)
	s.Mutex.Unlock()

	if err != nil {
		log.Printf("Error adding player to room %s: %v", roomID, err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Room is full"))
		conn.Close()
		return
	}

	err = conn.WriteJSON(map[string]string{"type": "client_id", "id": clientID, "role": player.Role})
	if err != nil {
		log.Printf("Error sending client ID to %s: %v", r.RemoteAddr, err)
		return
	}

	go s.broadcastGameState(s.Rooms[roomID])

	defer func() {
		log.Printf("Client disconnected from %s", r.RemoteAddr)
		s.Mutex.Lock()
		delete(s.Rooms[roomID].Players, player.ID)
		log.Printf("Player removed from room %s: %v", roomID, player)
		s.Mutex.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", r.RemoteAddr, err)
			break
		}
		var input map[string]interface{}
		if err := json.Unmarshal(msg, &input); err == nil {
			switch input["type"] {
			case "move":
				log.Printf("Received message from %s: %s", r.RemoteAddr, string(msg))
				s.handlePlayerMove(roomID, player.ID, input["direction"].(string))
			case "ping":
				// Responder inmediatamente con un pong
				timestamp := input["timestamp"]
				conn.WriteJSON(map[string]interface{}{
					"type":              "pong",
					"originalTimestamp": timestamp,
				})
			default:
				log.Printf("Received message from %s: %s", r.RemoteAddr, string(msg))
			}
		}
	}
}

func (s *Server) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	status := map[string]interface{}{
		"rooms": s.Rooms,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Println("Error encoding server status:", err)
		http.Error(w, "Could not encode server status", http.StatusInternalServerError)
	}
}
