package room

import "sync"

type Room struct {
	ID         string
	Players    [2]*Player
	SecretWord string
	State      string
	mu         sync.Mutex
}

type Player struct {
}
