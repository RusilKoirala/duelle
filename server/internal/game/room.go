package game

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

type RoomState string

const (
	Waiting  RoomState = "waiting"
	Playing  RoomState = "playing"
	Finished RoomState = "finished"
)

type Room struct {
	ID         string
	Players    []*Player
	SecretWord string
	State      RoomState
	mu         sync.RWMutex
}

// create a roomm
func NewRoom(id string, secretWord string) *Room {
	return &Room{
		ID:         id,
		Players:    make([]*Player, 0, 2),
		SecretWord: secretWord,
		State:      Waiting,
	}
}

// add someone
func (r *Room) AddPlayer(player *Player) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Players) >= 2 {
		return false
	}
	r.Players = append(r.Players, player)

	if len(r.Players) == 2 {
		r.State = Playing
	}
	return true
}

// remove them
func (r *Room) RemovePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.Players {
		if p.ID == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			break
		}
	}
	if len(r.Players) == 0 {
		r.State = Finished
	}
}

// get opponenttt
func (r *Room) GetOpponent(playerID string) *Player {
	r.mu.RUnlock()
	defer r.mu.Unlock()

	for _, p := range r.Players {
		if p.ID != playerID {
			return p
		}
	}
	return nil
}

// check if its full
func (r *Room) IsFull() bool {
	r.mu.Lock()
	defer r.mu.RUnlock()
	return len(r.Players) >= 2
}

// broadcast it
func (r *Room) Broadcast(ctx context.Context, message interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return
	}

	for _, player := range r.Players {
		if err := player.Send(ctx, data); err != nil {
			log.Printf("Send error to %s: %v", player.ID, err)
		}
	}
}
