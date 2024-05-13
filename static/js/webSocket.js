let socket = new WebSocket("ws://localhost:8080/websocket");

function updateGameUI(gameData) {
}

socket.onmessage = function(event) {
    let gameData = JSON.parse(event.data);
    updateGameUI(gameData);
};

function sendWebSocketMessage(message) {
    if (socket.readyState === WebSocket.OPEN) { 
        socket.send(JSON.stringify(message));
    } else {
        console.error("La connexion WebSocket n'est pas ouverte.");
    }
}
