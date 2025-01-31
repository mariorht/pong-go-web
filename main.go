package main

import (
    "log"
    "net/http"
    "pong-go-web/server" // Ensure this import path is correct
)

func main() {
    srv := server.NewServer()

    // Start the game loop
    go srv.StartGameLoop()

    // Servir archivos estáticos desde la carpeta "static/"
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    // Ruta WebSocket
    http.HandleFunc("/ws", srv.HandleWebSocket)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
