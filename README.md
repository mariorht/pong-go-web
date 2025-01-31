# Pong Go Web

This project is a web-based Pong game implemented using Go for the server and HTML5 Canvas for the client. The server handles WebSocket connections to manage real-time game state updates.

## Project Structure

```
pong-go-web/
│── main.go               # Entry point of the server
│── server/
│   ├── server.go         # Server initialization
│   ├── handlers.go       # WebSocket connection and event handling
│   ├── rooms.go          # Game room management
│   ├── game.go           # Game logic (physics, collisions, etc.)
│   ├── player.go         # Player structure and state
│── models/
│   ├── room.go           # Room structure definition
│   ├── player.go         # Player structure definition
│── static/
│   ├── index.html        # Web client
│   ├── game.js           # Frontend logic
│   ├── styles.css        # Client styles
│── go.mod                # Project dependencies
│── go.sum                # Dependency checksums
│── README.md             # Project overview and instructions
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- A modern web browser

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/mariorht/pong-go-web.git
    cd pong-go-web
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

### Running the Server

1. Start the server:

    ```sh
    go run main.go
    ```

2. Open your web browser and navigate to `http://localhost:8080`.

### Playing the Game

- The first player to connect will be assigned as Player 1 (left paddle).
- The second player to connect will be assigned as Player 2 (right paddle).
- Use the arrow keys to move your paddle up and down.
- The game will keep track of the score and reset the ball when a goal is scored.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
