package game

import (
	"context"

	"nhooyr.io/websocket"
)

type Player struct {
	ID      string
	Conn    *websocket.Conn
	Guesses []string
	Won     bool
}

func NewPlayer(id string, conn *websocket.Conn) *Player {
	return &Player{
		ID:      id,
		Conn:    conn,
		Guesses: make([]string, 0),
		Won:     false,
	}
}

func (p *Player) Send(ctx context.Context, data []byte) error {
	return p.Conn.Write(ctx, websocket.MessageText, data)
}
