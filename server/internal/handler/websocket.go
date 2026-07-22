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
}

func NewWSHandler(wordService *words.WordService) *WSHandler {
	return &WSHandler{
		wordService: wordService,
	}
}

type ClientMessage struct {
	Type string `json:"type"`
	Word string `json:"word"`
}

type ServerMessage struct {
	Type    string              `json:"type"`
	Valid   bool                `json:"valid,omitempty"`
	Message string              `json:"message,omitempty"`
	Results []game.LetterStatus `json:"results,omitempty"`
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

	log.Printf("WebSocket connected | Room: %s", roomID)

	ctx := context.Background()

	// connection alive
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := conn.Ping(ctx); err != nil {
				log.Printf("Ping failed: %v", err)
				return
			}
		}
	}()

	testSecret := "HELLO"

	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			log.Printf("Connection closed: %v", err)
			break
		}

		var msg ClientMessage

		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Invalid Message: %v", err)
		}

		log.Printf("Received: %s", string(data))

		if msg.Type == "guess" {

			if !h.wordService.IsValid(msg.Word) {
				response := ServerMessage{
					Type:    "error",
					Valid:   false,
					Message: "Not a valid word",
				}
				h.sendMessage(ctx, conn, response)
				continue
			}
			result := game.CheckGuess(testSecret, msg.Word)

			response := ServerMessage{
				Type:    "guess_result",
				Valid:   true,
				Results: result.Results,
			}

			h.sendMessage(ctx, conn, response)
		}
	}

	conn.Close(websocket.StatusNormalClosure, "Goodbye")
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
