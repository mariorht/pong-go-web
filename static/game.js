let playerRole = null;
let lastFrameTime = performance.now();
let frameCount = 0;
let fps = 0;
let lastPingTime = 0;

const ws = new WebSocket(`ws://${window.location.host}/ws`);

ws.onopen = () => {
    const statusEl = document.getElementById("status");
    statusEl.innerText = "Connected";
    statusEl.className = "status-item connected";
    console.log("Connected to server");
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
            console.log("Received client ID:", data.id);
            playerRole = data.role;
            document.getElementById("client-id").innerText = `Player: ${data.role}`;
        } else if (data.type === "pong") {
            const now = performance.now() * 1000;
            const rtt = now - data.originalTimestamp;
            document.getElementById("latency").innerText = `RTT: ${(rtt/1000).toFixed(2)} ms`;
        } else if (data.type === "game_state") {
            updateGame(data.state);
        }
    } catch (error) {
        console.error("Error parsing WebSocket message:", error);
    }
};

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
    const minutes = Math.floor(state.gameTime / 60);
    const seconds = state.gameTime % 60;
    ctx.fillText(`Time: ${minutes}:${seconds.toString().padStart(2, '0')}`, canvas.width / 2, 30);

    // Draw scores
    ctx.font = "30px Arial";
    ctx.fillText(`${state.score1} - ${state.score2}`, canvas.width / 2, 70);

    // Draw paddles
    ctx.fillStyle = "white";
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
        ws.send(JSON.stringify({ type: "move", direction: event.key, role: playerRole, sendTime }));
    }
});
