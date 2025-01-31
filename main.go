package main

import (
    "log"
    "net/http"
    "pong-server/server"
)

func main() {
    srv := server.NewServer()

    // Servir archivos estáticos desde la carpeta "static/"
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    // Ruta WebSocket
    http.HandleFunc("/ws", srv.HandleWebSocket)

    // Ruta de debug
    http.HandleFunc("/debug", srv.GetServerStatus)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
