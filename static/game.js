let gameConfig = {
    fieldWidth: 800,    // valores por defecto
    fieldHeight: 400,
    paddleWidth: 10,
    paddleHeight: 100,
    ballRadius: 10
};

let playerRole = null;
let lastFrameTime = performance.now();
let frameCount = 0;
let fps = 0;
let lastPingTime = 0;

let playerName = '';
let wsConnection = null;

function connectToGame() {
    if (playerName.trim() === '') {
        alert('Por favor, introduce tu nombre');
        return;
    }

    document.getElementById('find-game').disabled = true;
    document.getElementById('player-name').disabled = true;
    document.getElementById('status').innerText = 'Conectando...';

    wsConnection = new WebSocket(`ws://${window.location.host}/ws`);
    setupWebSocket(wsConnection);
}

// Mover la lógica de WebSocket a una función separada
function setupWebSocket(ws) {
    ws.onopen = () => {
        const statusEl = document.getElementById("status");
        statusEl.innerText = "Conectado";
        statusEl.className = "status-item connected";
        
        // Ocultar lobby inicial y mostrar lobby de espera
        document.getElementById('initial-lobby').style.display = 'none';
        document.getElementById('lobby').style.display = 'block';
        
        // Enviar nombre del jugador
        ws.send(JSON.stringify({ 
            type: "set_name",
            name: playerName 
        }));
    };

    ws.onerror = (error) => {
        console.error("WebSocket Error:", error);
        const statusEl = document.getElementById("status");
        statusEl.innerText = "Connection Error";
        statusEl.className = "status-item error";
    };

    ws.onclose = () => {
        const statusEl = document.getElementById("status");
        statusEl.innerText = "Disconnected";
        statusEl.className = "status-item error";
        console.log("Disconnected from server");
    };

    // Enviar ping cada segundo
    setInterval(() => {
        lastPingTime = performance.now() * 1000;
        ws.send(JSON.stringify({ type: "ping", timestamp: lastPingTime }));
    }, 1000);

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);

            if (data.type === "client_id") {
                playerRole = data.role;
                playerName = data.name;
                const roleText = data.role === "player1" ? "Jugador 1" : "Jugador 2";
                document.getElementById("client-id").innerText = `${roleText}`;
                document.getElementById("status").innerText = data.name;
            } else if (data.type === "pong") {
                const now = performance.now();
                const rtt = now - data.originalTimestamp / 1000;
                document.getElementById("latency").innerText = `RTT: ${rtt.toFixed(2)} ms`;
            } else if (data.type === "game_state") {
                if (data.config) {
                    gameConfig = data.config;
                    canvas.width = gameConfig.fieldWidth;
                    canvas.height = gameConfig.fieldHeight;
                }

                const lobby = document.getElementById("lobby");
                const gameCanvas = document.getElementById("gameCanvas");
                const countdown = document.getElementById("countdown");
                
                // console.log("Game state:", data.state.state); // Debug

                switch(data.state.state) {
                    case "waiting":
                        lobby.style.display = "block";
                        gameCanvas.style.display = "none";
                        lobby.querySelector("h2").textContent = "Waiting for opponent...";
                        countdown.textContent = "";
                        break;
                        
                    case "starting":
                        lobby.style.display = "block";
                        gameCanvas.style.display = "none";
                        const startTime = new Date(data.state.startTime);
                        const timeLeft = Math.ceil((startTime - new Date()) / 1000);
                        countdown.textContent = `Game starts in ${timeLeft} seconds...`;
                        break;
                        
                    case "playing":
                        lobby.style.display = "none";
                        gameCanvas.style.display = "block";
                        updateGame(data.state);
                        break;

                    case "finished":
                        lobby.style.display = "block";
                        gameCanvas.style.display = "none";
                        // console.log("Game finished. Winner:", data.state.winner, "My role:", playerRole); // Debug
                        const isWinner = data.state.winner === playerRole;
                        let resultMessage;
                        if (isWinner) {
                            resultMessage = `¡Has ganado! (${data.state.score1} - ${data.state.score2})`;
                        } else {
                            resultMessage = `Has perdido (${data.state.score1} - ${data.state.score2})`;
                        }
                        
                        if (data.state.winReason === "opponent_disconnected") {
                            resultMessage += " - Victoria por abandono";
                        }
                        
                        lobby.querySelector("h2").textContent = resultMessage;
                        countdown.textContent = "";
                        break;
                }
            }
            // ...rest of the message handling...
        } catch (error) {
            console.error("Error parsing WebSocket message:", error);
        }
    };
}

const canvas = document.getElementById("gameCanvas");
const ctx = canvas.getContext("2d");

function updateGame(state) {
    const now = Date.now();
    frameCount++;
    if (now - lastFrameTime >= 1000) {
        fps = frameCount;
        frameCount = 0;
        lastFrameTime = now;
        document.getElementById("fps").innerText = `FPS: ${fps}`;
    }

    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // Draw game time
    ctx.font = "20px Arial";
    ctx.textAlign = "center";
    ctx.fillStyle = "white";
    const minutes = Math.floor(state.gameTime / 60);
    const seconds = state.gameTime % 60;
    ctx.fillText(`Time: ${minutes}:${seconds.toString().padStart(2, '0')}`, canvas.width / 2, 30);

    // Draw scores
    ctx.font = "30px Arial";
    ctx.fillText(`${state.score1} - ${state.score2}`, canvas.width / 2, 70);

    // Draw paddles
    ctx.fillRect(state.paddle1.x, state.paddle1.y, state.paddle1.width, state.paddle1.height);
    ctx.fillRect(state.paddle2.x, state.paddle2.y, state.paddle2.width, state.paddle2.height);

    // Draw balls
    state.balls.forEach(ball => {
        ctx.beginPath();
        ctx.arc(ball.x, ball.y, ball.radius, 0, Math.PI * 2);
        ctx.fill();
    });
}

document.addEventListener("keydown", (event) => {
    if (event.key === "ArrowUp" || event.key === "ArrowDown") {
        const sendTime = Math.floor(performance.now() * 1000); // Microseconds
        wsConnection.send(JSON.stringify({ type: "move", direction: event.key, role: playerRole, sendTime }));
    }
});

// Inicialización
document.addEventListener('DOMContentLoaded', () => {
    const nameInput = document.getElementById('player-name');
    const findGameBtn = document.getElementById('find-game');

    nameInput.addEventListener('input', (e) => {
        playerName = e.target.value;
        findGameBtn.disabled = playerName.trim() === '';
    });

    findGameBtn.addEventListener('click', connectToGame);
});
