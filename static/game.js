const ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = () => {
    document.getElementById("status").innerText = "Connected to server";
    console.log("Connected to server");
};

ws.onerror = (error) => {
    console.error("WebSocket Error:", error);
    document.getElementById("status").innerText = "Error connecting to server";
};

ws.onclose = () => {
    document.getElementById("status").innerText = "Disconnected";
    console.log("Disconnected from server");
};

ws.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);
        if (data.type === "client_id") {
            console.log("Received client ID:", data.id);
            document.getElementById("client-id").innerText = `Client ID: ${data.id}`;
        } else if (data.type === "game_state") {
            console.log("Received game state:", data.state);
            updateGame(data.state);
        }
    } catch (error) {
        console.error("Error parsing WebSocket message:", error);
    }
};

const canvas = document.getElementById("gameCanvas");
const ctx = canvas.getContext("2d");

function updateGame(state) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    // Draw paddles and ball based on the state received from the server
    ctx.fillStyle = "white";
    ctx.fillRect(state.paddle1.x, state.paddle1.y, state.paddle1.width, state.paddle1.height);
    ctx.fillRect(state.paddle2.x, state.paddle2.y, state.paddle2.width, state.paddle2.height);
    ctx.beginPath();
    ctx.arc(state.ball.x, state.ball.y, state.ball.radius, 0, Math.PI * 2);
    ctx.fill();
}

document.addEventListener("keydown", (event) => {
    if (event.key === "ArrowUp" || event.key === "ArrowDown") {
        ws.send(JSON.stringify({ type: "move", direction: event.key }));
    }
});
