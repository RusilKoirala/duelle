package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rusilkoirala/duelle/internal/game"
	"github.com/rusilkoirala/duelle/internal/words"
	"nhooyr.io/websocket"
)

type WSHandler struct {
	wordService *words.WordService
	manager     *game.Manager
}

func NewWSHandler(wordService *words.WordService, manager *game.Manager) *WSHandler {
	return &WSHandler{
		wordService: wordService,
		manager:     manager,
	}
}

type ClientMessage struct {
	Type string `json:"type"`
	Word string `json:"word"`
}

type ServerMessage struct {
	Type            string              `json:"type"`
	Valid           bool                `json:"valid,omitempty"`
	Message         string              `json:"message,omitempty"`
	Results         []game.LetterStatus `json:"results,omitempty"`
	OpponentGuesses int                 `json:"opponent_guesses,omitempty"`
	RoomState       string              `json:"room_state,omitempty"`
	PlayerCount     int                 `json:"player_count,omitempty"`
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Missing room parameter", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	playerID := generatePlayerID()

	room := h.manager.GetOrCreateRoom(roomID, h.wordService.GetRandomWord())
	player := game.NewPlayer(playerID, conn)

	if !room.AddPlayer(player) {
		h.sendError(context.Background(), conn, "Room is full")
		conn.Close(websocket.StatusNormalClosure, "Room full")
		return
	}

	h.sendRoomState(context.Background(), conn, room)

	ctx := context.Background()

	// connection alivee
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := conn.Ping(ctx); err != nil {
				return
			}
		}
	}()

	h.handleMessages(ctx, room, player)

	room.RemovePlayer(playerID)
	conn.Close(websocket.StatusNormalClosure, "Goodbye")
}

func (h *WSHandler) handleMessages(ctx context.Context, room *game.Room, player *game.Player) {
	for {
		_, data, err := player.Conn.Read(ctx)
		if err != nil {
			log.Printf("Player %s disconnected", player.ID)
			break
		}

		var msg ClientMessage

		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		if msg.Type == "guess" {
			h.handleGuess(ctx, room, player, msg.Word)
		}

	}
}

func (h *WSHandler) handleGuess(ctx context.Context, room *game.Room, player *game.Player, word string) {
	if !h.wordService.IsValid(word) {
		h.sendError(ctx, player.Conn, "Not a valid word")
		return
	}

	result := game.CheckGuess(room.SecretWord, word)
	player.Guesses = append(player.Guesses, word)

	opponent := room.GetOpponent(player.ID)

	opponentGuesses := 0

	if opponent != nil {
		opponentGuesses = len(opponent.Guesses)
	}

	response := ServerMessage{
		Type:            "guess_result",
		Valid:           true,
		Results:         result.Results,
		OpponentGuesses: opponentGuesses,
	}

	responseData, _ := json.Marshal(response)
	player.Send(ctx, responseData)

	if allCorrect(result.Results) {
		player.Won = true
		room.State = game.Finished
	}

	if opponent != nil {
		opponentMsg := ServerMessage{
			Type:            "opponent_guessed",
			OpponentGuesses: len(player.Guesses),
		}

		opponentData, _ := json.Marshal(opponentMsg)
		opponent.Send(ctx, opponentData)
	}

}

func (h *WSHandler) sendMessage(ctx context.Context, conn *websocket.Conn, msg ServerMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return
	}

	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		log.Printf("Write error: %v", err)
	}
}

func (h *WSHandler) sendRoomState(ctx context.Context, conn *websocket.Conn, room *game.Room) {
	msg := ServerMessage{
		Type:        "room_state",
		RoomState:   string(room.State),
		PlayerCount: len(room.Players),
	}
	data, _ := json.Marshal(msg)
	conn.Write(ctx, websocket.MessageText, data)
}

func generatePlayerID() string {
	return "P" + time.Now().Format("150405")
}
func allCorrect(results []game.LetterStatus) bool {
	for _, r := range results {
		if r != game.Correct {
			return false
		}
	}
	return true
}

func (h *WSHandler) sendError(ctx context.Context, conn *websocket.Conn, message string) {
	msg := ServerMessage{
		Type:    "error",
		Message: message,
	}
	data, _ := json.Marshal(msg)
	conn.Write(ctx, websocket.MessageText, data)
}
