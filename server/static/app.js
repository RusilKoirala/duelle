class DuelleGame {
    constructor() {
        this.roomId = null;
        this.ws = null;
        this.currentRow = 0;
        this.currentCol = 0;
        this.currentGuess = '';
        this.guesses = [];
        this.opponentGuesses = 0;
        this.gameActive = false;
        this.startTime = null;
        this.timerInterval = null;
        this.keyboardState = {};
        this.stats = this.loadStats();

        this.sounds = {
            keypress: this.createSound(200, 0.1, 'sine'),
            invalid: this.createSound(100, 0.2, 'sawtooth'),
            win: this.createSound(440, 0.3, 'sine')
        };

        this.checkURLRoom();
        this.setupMenu();
        this.setupStatsModal();
    }

    createSound(frequency, duration, type) {
        return () => {
            const audioContext = new (window.AudioContext || window.webkitAudioContext)();
            const oscillator = audioContext.createOscillator();
            const gainNode = audioContext.createGain();

            oscillator.connect(gainNode);
            gainNode.connect(audioContext.destination);

            oscillator.frequency.value = frequency;
            oscillator.type = type;
            gainNode.gain.setValueAtTime(0.1, audioContext.currentTime);
            gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + duration);

            oscillator.start(audioContext.currentTime);
            oscillator.stop(audioContext.currentTime + duration);
        };
    }

    loadStats() {
        const saved = localStorage.getItem('duelle_stats');
        return saved ? JSON.parse(saved) : {
            gamesPlayed: 0,
            gamesWon: 0,
            currentStreak: 0,
            totalGuesses: 0
        };
    }

    saveStats() {
        localStorage.setItem('duelle_stats', JSON.stringify(this.stats));
    }

    updateStats(won, guesses) {
        this.stats.gamesPlayed++;
        if (won) {
            this.stats.gamesWon++;
            this.stats.currentStreak++;
            this.stats.totalGuesses += guesses;
        } else {
            this.stats.currentStreak = 0;
        }
        this.saveStats();
    }

    setupStatsModal() {
        const statsBtn = document.getElementById('stats-btn');
        const modal = document.getElementById('stats-modal');
        const closeBtn = modal.querySelector('.close');

        statsBtn.addEventListener('click', () => {
            this.showStats();
        });

        closeBtn.addEventListener('click', () => {
            modal.classList.add('hidden');
        });

        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.classList.add('hidden');
            }
        });
    }

    showStats() {
        const winPercent = this.stats.gamesPlayed > 0
            ? Math.round((this.stats.gamesWon / this.stats.gamesPlayed) * 100)
            : 0;
        const avgGuesses = this.stats.gamesWon > 0
            ? (this.stats.totalGuesses / this.stats.gamesWon).toFixed(1)
            : 0;

        document.getElementById('games-played').textContent = this.stats.gamesPlayed;
        document.getElementById('win-percent').textContent = winPercent;
        document.getElementById('current-streak').textContent = this.stats.currentStreak;
        document.getElementById('avg-guesses').textContent = avgGuesses;

        document.getElementById('stats-modal').classList.remove('hidden');
    }

    checkURLRoom() {
        const params = new URLSearchParams(window.location.search);
        const roomFromURL = params.get('room');
        if (roomFromURL && roomFromURL.length === 6) {
            this.roomId = roomFromURL.toUpperCase();
            this.startGame();
        }
    }

    setupMenu() {
        document.getElementById('create-room-btn').addEventListener('click', () => {
            this.createRoom();
        });

        document.getElementById('join-room-btn').addEventListener('click', () => {
            document.getElementById('join-input-container').classList.remove('hidden');
        });

        document.getElementById('join-submit-btn').addEventListener('click', () => {
            const code = document.getElementById('room-code-input').value.trim().toUpperCase();
            if (code.length === 6) {
                this.joinRoom(code);
            }
        });

        document.getElementById('room-code-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                const code = document.getElementById('room-code-input').value.trim().toUpperCase();
                if (code.length === 6) {
                    this.joinRoom(code);
                }
            }
        });
    }

    createRoom() {
        this.roomId = Math.random().toString(36).substr(2, 6).toUpperCase();
        this.startGame();
    }

    joinRoom(code) {
        this.roomId = code;
        this.startGame();
    }

    startGame() {
        document.getElementById('menu').classList.add('hidden');
        document.getElementById('game').classList.remove('hidden');

        this.initUI();
        this.initKeyboard();
        this.connectWebSocket();
        this.setupKeyboard();
        this.updateURL();
    }

    updateURL() {
        const url = new URL(window.location);
        url.searchParams.set('room', this.roomId);
        window.history.replaceState({}, '', url);
    }

    initUI() {
        const board = document.getElementById('board');
        board.innerHTML = '';

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

            board.appendChild(rowDiv);
        }

        document.getElementById('room-code').textContent = `Room: ${this.roomId}`;
        this.updateOpponentInfo();
    }

    initKeyboard() {
        const keys = document.querySelectorAll('.key');
        keys.forEach(key => {
            key.addEventListener('click', () => {
                const keyValue = key.dataset.key;
                if (keyValue === 'Enter') {
                    this.submitGuess();
                } else if (keyValue === 'Backspace') {
                    this.deleteLetter();
                } else {
                    this.addLetter(keyValue);
                }
            });
        });
    }

    updateKeyboard(letter, status) {
        const key = document.querySelector(`[data-key="${letter}"]`);
        if (!key) return;

        const currentStatus = this.keyboardState[letter];
        if (currentStatus === 'correct') return;
        if (currentStatus === 'present' && status === 'absent') return;

        this.keyboardState[letter] = status;
        key.classList.remove('correct', 'present', 'absent');
        key.classList.add(status);
    }

    startTimer() {
        if (!this.startTime) {
            this.startTime = Date.now();
        }

        this.timerInterval = setInterval(() => {
            const elapsed = Math.floor((Date.now() - this.startTime) / 1000);
            const minutes = Math.floor(elapsed / 60);
            const seconds = elapsed % 60;
            document.getElementById('timer').textContent =
                `${minutes}:${seconds.toString().padStart(2, '0')}`;
        }, 1000);
    }

    stopTimer() {
        if (this.timerInterval) {
            clearInterval(this.timerInterval);
        }
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?room=${this.roomId}`;

        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
            this.showMessage('Waiting for opponent...', '');
        };

        this.ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            this.handleServerMessage(msg);
        };

        this.ws.onerror = () => {
            this.showMessage('Connection error', 'error');
        };

        this.ws.onclose = () => {
            this.showMessage('Opponent disconnected', 'error');
            this.gameActive = false;
            this.stopTimer();
        };
    }

    handleServerMessage(msg) {
        if (msg.type === 'error') {
            this.showMessage(msg.message, 'error');
            this.shakeRow();
            this.sounds.invalid();
            this.currentGuess = '';
            this.currentCol = 0;
            this.updateCurrentRow();
            return;
        }

        if (msg.type === 'room_state') {
            if (msg.room_state === 'waiting') {
                this.showMessage('Waiting for opponent...', '');
                this.gameActive = false;
            } else if (msg.room_state === 'playing') {
                this.showMessage('Game started!', 'success');
                this.gameActive = true;
                this.startTimer();
                setTimeout(() => this.showMessage('', ''), 2000);
            }
            return;
        }

        if (msg.type === 'guess_result') {
            this.displayGuess(this.guesses[this.guesses.length - 1], msg.results);
            this.opponentGuesses = msg.opponent_guesses;
            this.updateOpponentInfo();
            return;
        }

        if (msg.type === 'opponent_guessed') {
            this.opponentGuesses = msg.opponent_guesses;
            this.updateOpponentInfo();
            return;
        }

        if (msg.type === 'game_over') {
            this.handleGameOver(msg);
            return;
        }
    }

    handleGameOver(msg) {
        this.gameActive = false;
        this.stopTimer();

        const won = msg.winner === 'you';
        this.updateStats(won, this.currentRow + 1);

        if (won) {
            this.showMessage('You won!', 'success');
            this.sounds.win();
            this.celebrateWin();
        } else if (msg.winner === 'opponent') {
            this.showMessage('Opponent won!', 'error');
        } else {
            this.showMessage('Game over', 'error');
        }
    }

    celebrateWin() {
        for (let i = 0; i < 5; i++) {
            const tile = document.querySelector(`[data-row="${this.currentRow - 1}"][data-col="${i}"]`);
            setTimeout(() => {
                tile.classList.add('bounce');
            }, i * 100);
        }
    }

    shakeRow() {
        const row = document.querySelectorAll(`[data-row="${this.currentRow}"]`);
        row.forEach(tile => {
            tile.classList.add('shake');
            setTimeout(() => tile.classList.remove('shake'), 500);
        });
    }

    updateOpponentInfo() {
        document.getElementById('opponent-info').textContent =
            `Opponent: ${this.opponentGuesses}/6`;
    }

    setupKeyboard() {
        document.addEventListener('keydown', (e) => {
            if (!this.gameActive) return;
            if (this.currentRow >= 6) return;

            if (e.key === 'Enter') {
                this.submitGuess();
            } else if (e.key === 'Backspace') {
                this.deleteLetter();
            } else if (e.key.length === 1 && e.key.match(/[a-z]/i)) {
                this.addLetter(e.key.toUpperCase());
            }
        });
    }

    addLetter(letter) {
        if (this.currentCol < 5) {
            this.currentGuess += letter;
            this.updateCurrentRow();
            this.currentCol++;
            this.sounds.keypress();
        }
    }

    deleteLetter() {
        if (this.currentCol > 0) {
            this.currentGuess = this.currentGuess.slice(0, -1);
            this.currentCol--;
            this.updateCurrentRow();
            this.sounds.keypress();
        }
    }

    updateCurrentRow() {
        for (let i = 0; i < 5; i++) {
            const tile = document.querySelector(`[data-row="${this.currentRow}"][data-col="${i}"]`);
            if (i < this.currentGuess.length) {
                tile.textContent = this.currentGuess[i];
                tile.classList.add('filled');
            } else {
                tile.textContent = '';
                tile.classList.remove('filled');
            }
        }
    }

    submitGuess() {
        if (this.currentGuess.length !== 5) {
            this.showMessage('Not enough letters', 'error');
            this.shakeRow();
            this.sounds.invalid();
            return;
        }

        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            this.showMessage('Not connected', 'error');
            return;
        }

        this.ws.send(JSON.stringify({
            type: "guess",
            word: this.currentGuess
        }));

        this.guesses.push(this.currentGuess);
        this.showMessage('', '');
    }

    displayGuess(guess, results) {
        for (let i = 0; i < 5; i++) {
            const tile = document.querySelector(`[data-row="${this.currentRow}"][data-col="${i}"]`);

            setTimeout(() => {
                tile.classList.add('flip');
                setTimeout(() => {
                    tile.classList.add(results[i]);
                    this.updateKeyboard(guess[i], results[i]);
                }, 250);
            }, 150 * i);
        }

        if (results.every(r => r === 'correct')) {
            setTimeout(() => {
                this.showMessage('You won!', 'success');
                this.gameActive = false;
                this.stopTimer();
            }, 1000);
        } else {
            this.currentRow++;
            this.currentCol = 0;
            this.currentGuess = '';

            if (this.currentRow >= 6) {
                setTimeout(() => {
                    this.showMessage('Game over', 'error');
                    this.gameActive = false;
                    this.stopTimer();
                }, 1000);
            }
        }
    }

    showMessage(text, type) {
        const messageEl = document.getElementById('message');
        messageEl.textContent = text;
        messageEl.className = `message ${type}`;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new DuelleGame();
});
