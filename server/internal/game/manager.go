package game

import (
	"log"
	"sync"
)

type Manager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

// new manager
func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

// create a room
func (m *Manager) CreateRoom(roomID string, secretWord string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	room := NewRoom(roomID, secretWord)
	m.rooms[roomID] = room
	log.Printf("room created: %s", roomID)
	return room
}

// get room by id
func (m *Manager) GetRoom(roomID string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rooms[roomID]
}

// get existing room or create one
func (m *Manager) GetOrCreateRoom(roomID string, secretWord string) *Room {
	room := m.GetRoom(roomID)
	if room == nil {
		return m.CreateRoom(roomID, secretWord)
	}
	return room
}

// delete a room
func (m *Manager) DeleteRoom(roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rooms, roomID)
	log.Printf("room deleted: %s", roomID)
}
