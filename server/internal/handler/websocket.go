package handler

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

type WSHandler struct {
}

func NewWSHandler() *WSHandler {
	return &WSHandler{}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Missing room parameter", http.StatusBadRequest)
		return
	}

	// websocket timeee
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Printf("WebSocket upgrade failed: %d ", err)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "Connection closed")

	log.Printf("websocket connected | Room: %s", roomID)

	//ctx := context.Background()
	conn.Close(websocket.StatusNormalClosure, "Echo test")
}
