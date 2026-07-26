package game

import (
	"log"
	"sync"
)

type Manager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

// create new manager
func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

// create room
func (m *Manager) CreateRoom(roomID string, secretWord string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	room := NewRoom(roomID, secretWord)
	m.rooms[roomID] = room
	log.Printf("Room created: %s", roomID)
	return room
}

// get the room
func (m *Manager) GetRoom(roomID string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rooms[roomID]
}

// get the room if not create
func (m *Manager) GetOrCreateRoom(roomID string, secretWord string) *Room {
	room := m.GetRoom(roomID)
	if room == nil {
		return m.CreateRoom(roomID, secretWord)
	}
	return room
}

// delete room
func (m *Manager) DeleteRoom(roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rooms, roomID)
	log.Printf("Room deleted: %s", roomID)
}
