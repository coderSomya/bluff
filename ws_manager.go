// ws_manager.go
package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	mu           sync.RWMutex
	clients      map[string][]*websocket.Conn  // gameID -> list of connections
	playerConns  map[*websocket.Conn]string    // connection -> playerID
	connPlayers  map[string]*websocket.Conn    // playerID -> connection
	connGames    map[*websocket.Conn]string    // connection -> gameID
}

var manager = ClientManager{
	clients:     make(map[string][]*websocket.Conn),
	playerConns: make(map[*websocket.Conn]string),
	connPlayers: make(map[string]*websocket.Conn),
	connGames:   make(map[*websocket.Conn]string),
}

// Add a client to a game with player mapping
func (cm *ClientManager) AddClient(gameID string, conn *websocket.Conn, playerID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// Add to game's connection list
	cm.clients[gameID] = append(cm.clients[gameID], conn)
	
	// Create bidirectional mappings
	cm.playerConns[conn] = playerID
	cm.connPlayers[playerID] = conn
	cm.connGames[conn] = gameID
}

// Get player ID for a connection
func (cm *ClientManager) GetPlayerID(conn *websocket.Conn) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.playerConns[conn]
}

// Get connection for a player
func (cm *ClientManager) GetPlayerConnection(playerID string) *websocket.Conn {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.connPlayers[playerID]
}

// Remove a client when they disconnect
func (cm *ClientManager) RemoveClient(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// Get player and game info before removing
	playerID := cm.playerConns[conn]
	gameID := cm.connGames[conn]
	
	// Remove from game's connection list
	if conns, exists := cm.clients[gameID]; exists {
		for i, c := range conns {
			if c == conn {
				cm.clients[gameID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		// Clean up empty game
		if len(cm.clients[gameID]) == 0 {
			delete(cm.clients, gameID)
		}
	}
	
	// Remove all mappings
	delete(cm.playerConns, conn)
	delete(cm.connPlayers, playerID)
	delete(cm.connGames, conn)
}

// Send message to specific player
func (cm *ClientManager) SendToPlayer(playerID string, message []byte) error {
	cm.mu.RLock()
	conn := cm.connPlayers[playerID]
	cm.mu.RUnlock()
	
	if conn != nil {
		return conn.WriteMessage(websocket.TextMessage, message)
	}
	return fmt.Errorf("player %s not connected", playerID)
}

// Broadcast to all players in a game with personalized messages
func (cm *ClientManager) BroadcastPersonalized(gameID string, messageFunc func(playerID string) []byte) {
	cm.mu.RLock()
	conns := cm.clients[gameID]
	players := make([]string, len(conns))
	for i, conn := range conns {
		players[i] = cm.playerConns[conn]
	}
	cm.mu.RUnlock()
	
	for i, conn := range conns {
		if players[i] != "" {
			message := messageFunc(players[i])
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// Original broadcast method (keep for backward compatibility)
func (cm *ClientManager) Broadcast(gameID string, message []byte) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if conns, ok := cm.clients[gameID]; ok {
		for _, conn := range conns {
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}