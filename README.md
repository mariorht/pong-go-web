# Pong Go Web

[![Release](https://img.shields.io/github/v/release/mariorht/pong-go-web?style=flat-square)](https://github.com/mariodev/pong-go-web/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mariorht/pong-go-web?style=flat-square)](go.mod)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mariorht/pong-go-web/release.yml?style=flat-square)](https://github.com/mariorht/pong-go-web/actions/workflows/release.yml)
[![License](https://img.shields.io/github/license/mariorht/pong-go-web?style=flat-square)](LICENSE)

This project is a web-based Pong game implemented using Go for the server and HTML5 Canvas for the client. The server handles WebSocket connections to manage real-time game state updates.

## Current State

![Pong Game Screenshot](docs/images/game.png)

Current implementation shows a basic two-player pong game with real-time WebSocket communication.


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

- Docker and Docker Compose
OR
- Go 1.16 or later
- A modern web browser

### Running with Docker

1. Build and start the container:
    ```sh
    docker-compose up --build
    ```

2. Open your web browser and navigate to `http://localhost:8080`

3. To stop the server:
    ```sh
    docker-compose down
    ```

### Running without Docker

1. Install dependencies:
    ```sh
    go mod tidy
    ```

2. Start the server:
    ```sh
    go run main.go
    ```

3. Open your web browser and navigate to `http://localhost:8080`

### Playing the Game

- The first player to connect will be assigned as Player 1 (left paddle).
- The second player to connect will be assigned as Player 2 (right paddle).
- Use the arrow keys to move your paddle up and down.
- The game will keep track of the score and reset the ball when a goal is scored.

## TODO List

- [ ] Fix double collision issues when balls overlap
- [ ] Add paddle acceleration mechanics
- [ ] Implement support for multiple game rooms
- [ ] Add game mode selection and custom rules:
  - Victory by score limit
  - Victory by goal difference
  - Time-based matches
  - Multiple balls spawning
- [ ] Allow players to choose custom names
- [ ] Add test coverage:
  - Unit tests for game logic
  - Integration tests for WebSocket communication
  - E2E tests for full game scenarios
  - Test coverage report in CI

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
