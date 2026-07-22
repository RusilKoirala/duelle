package handler

import (
	"context"
	"log"
	"net/http"
	"time"

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

	// Upgrade to WebSocket
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	log.Printf("✅ WebSocket connected | Room: %s", roomID)

	ctx := context.Background()

	// Keep connection alive with ping
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := conn.Ping(ctx); err != nil {
				log.Printf("Ping failed: %v", err)
				return
			}
		}
	}()

	// Read messages loop - KEEP CONNECTION OPEN
	for {
		msgType, data, err := conn.Read(ctx)
		if err != nil {
			log.Printf("Connection closed: %v", err)
			break
		}

		log.Printf("📨 Received: %s", string(data))

		// Echo back (Phase 1)
		if err := conn.Write(ctx, msgType, data); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}

	// Clean close when loop exits
	conn.Close(websocket.StatusNormalClosure, "Goodbye")
}
