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
        const guess = input.value.trim();

        if (guess.length !== 5) {
            alert('Word must be 5 letters!');
            return;
        }

        console.log('Submitting guess:', guess);


     
        this.displayGuess(guess);

        input.value = '';
    }

    displayGuess(guess) {
        if (this.currentRow >= 6) {
            alert('No more guesses!');
            return;
        }

        for (let i = 0; i < 5; i++) {
            const tile = document.querySelector(`[data-row="${this.currentRow}"][data-col="${i}"]`);
            tile.textContent = guess[i];
            tile.classList.add('filled');

           
            setTimeout(() => {
                const colors = ['correct', 'present', 'absent'];
                const randomColor = colors[Math.floor(Math.random() * colors.length)];
                tile.classList.add(randomColor);
            }, 200 * i);
        }

        this.currentRow++;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    console.log('Duelle');
    new DuelleGame();
});
