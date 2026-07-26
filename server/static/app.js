class DuelleGame {
    constructor(roomId) {
        this.roomId = roomId || this.generateRoomId();
        this.ws = null;
        this.currentRow = 0;
        this.guesses = [];

        this.initUI();
        this.connectWebSocket();
        this.setupEventListeners();
    }

    generateRoomId() {
        return 'ROOM' + Math.random().toString(36).substr(2, 6).toUpperCase();
    }

    initUI() {
        const yourBoard = document.getElementById('your-board');

        for (let row = 0; row < 6; row++) {
            const rowDiv = document.createElement('div');
            rowDiv.className = 'row';

            for (let col = 0; col < 5; col++) {
                const tile = document.createElement('div');
                tile.className = 'tile';
                tile.dataset.row = row;
                tile.dataset.col = col;
                rowDiv.appendChild(tile);
            }

            yourBoard.appendChild(rowDiv);
        }
        document.getElementById('room-code').textContent = `Room: ${this.roomId}`;
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?room=${this.roomId}`;

        console.log('Connecting to:', wsUrl);

        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
            console.log("WebSocket connected");
            document.getElementById('connection-status').textContent = '🟢 Connected';
        };

        this.ws.onmessage = (event) => {
            console.log('Message received:', event.data);
            const msg = JSON.parse(event.data);
            this.handleServerMessage(msg);
        };

        this.ws.onerror = (error) => {
            console.error("WebSocket error", error);
            document.getElementById('connection-status').textContent = "🔴 Error";
        };

        this.ws.onclose = () => {
            console.log("WebSocket Closed");
            document.getElementById('connection-status').textContent = "⚫ Disconnected";
        };
    }

    handleServerMessage(msg) {
        if (msg.type === 'error') {
            alert(msg.message);
            return;
        }

        if (msg.type === 'guess_result') {
            this.displayGuess(this.guesses[this.guesses.length - 1], msg.results);
        }
    }

    setupEventListeners() {
        const input = document.getElementById('guess-input');
        const submitBtn = document.getElementById('submit-btn');

        input.addEventListener('input', (e) => {
            e.target.value = e.target.value.toUpperCase();
        });

        input.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.submitGuess();
            }
        });

        submitBtn.addEventListener('click', () => {
            this.submitGuess();
        });
    }

    submitGuess() {
        const input = document.getElementById('guess-input');
        const guess = input.value.trim().toUpperCase();

        if (guess.length !== 5) {
            alert('Word must be 5 letters!');
            return;
        }

        if (this.currentRow >= 6) {
            alert("No more guesses!");
            return;
        }

        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            alert("Not connected to server");
            return;
        }

        console.log('Submitting guess:', guess);

        this.ws.send(JSON.stringify({
            type: "guess",
            word: guess
        }));

        this.guesses.push(guess);
        input.value = "";
    }

    displayGuess(guess, results) {
        for (let i = 0; i < 5; i++) {
            const tile = document.querySelector(`[data-row="${this.currentRow}"][data-col="${i}"]`);
            tile.textContent = guess[i];
            tile.classList.add('filled');

            setTimeout(() => {
                tile.classList.add(results[i]);
            }, 200 * i);
        }

        this.currentRow++;

        if (results.every(r => r === 'correct')) {
            setTimeout(() => {
                alert('🎉 You won!');
            }, 1500);
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    console.log('🎮 Duelle Phase 2');
    new DuelleGame();
});
