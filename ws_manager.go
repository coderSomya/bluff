// ws_manager.go
package main

import (
	"sync"
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	mu       sync.RWMutex //mutex to prevent race conditions
	clients  map[string][]*websocket.Conn // gameID -> list of connections
}

var manager = ClientManager{
	clients: make(map[string][]*websocket.Conn),
}

// Add a client to a game
func (cm *ClientManager) AddClient(gameID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[gameID] = append(cm.clients[gameID], conn)
}

// Broadcast a message to all clients in a game
func (cm *ClientManager) Broadcast(gameID string, message []byte) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if conns, ok := cm.clients[gameID]; ok {
		for _, conn := range conns {
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
