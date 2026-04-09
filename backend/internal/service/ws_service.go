package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WSManager struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewWSManager() *WSManager {
	return &WSManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (m *WSManager) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Non-websocket requests to /ws are expected sometimes (e.g. opening in browser).
		// We return early to avoid storing nil connections.
		log.Printf("websocket upgrade failed: %v", err)
		return
	}

	m.mu.Lock()
	m.clients[conn] = true
	m.mu.Unlock()
}

func (m *WSManager) Broadcast(data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for c := range m.clients {
		if c == nil {
			delete(m.clients, c)
			continue
		}
		if err := c.WriteJSON(data); err != nil {
			_ = c.Close()
			delete(m.clients, c)
		}
	}
}