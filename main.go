package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/coderSomya/bluff/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origin
		return true
	},
	// Add buffer sizes to prevent issues
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request:", r.Method, r.URL.Path)

	//cors ki bakchodi
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var gameID string

	respondJSON(conn, "connected", map[string]string{"status": "connected"});

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket unexpected close error: %v", err)
			} else {
				fmt.Printf("WebSocket read error: %v\n", err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			fmt.Println("Invalid JSON:", err)
			continue
		}

		switch wsMsg.Type {

		case "newGame":
			var payload NewGamePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				fmt.Println("Invalid newGame payload:", err)
				continue
			}
			gameID = "123"
			newGame := createGame(payload.Creator, gameID)
			respondJSON(conn, "newGame", newGame)

			manager.AddClient(gameID, conn)

		case "newPlayer":
			var payload NewPlayerPayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				fmt.Println("Invalid newPlayer payload:", err)
				continue
			}
			
			// Check if game exists before adding player
			game, _ := GetGameByID(payload.GameID)
			if game == nil {
				fmt.Printf("Game %s not found\n", payload.GameID)
				respondJSON(conn, "error", map[string]string{
					"message": "Game not found",
				})
				continue
			}
			
			// Add player to game
			if _,err := addPlayerToGame(payload.GameID, payload.Player); err != nil {
				fmt.Printf("Failed to add player to game: %v\n", err)
				respondJSON(conn, "error", map[string]string{
					"message": "Failed to join game",
				})
				continue
			}
			
			gameID = payload.GameID
			manager.AddClient(gameID, conn)
			
			// Send updated game state to the new player
			updatedGame, _ := GetGameByID(gameID)
			respondJSON(conn, "gameJoined", updatedGame)
			
			// Broadcast to other players that someone joined
			manager.Broadcast(gameID, mustMarshal("playerJoined", map[string]interface{}{
				"player": payload.Player,
				"game":   updatedGame,
			}))

		case "startGame":
			var payload StartGamePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}
			game := startGame(payload.GameID)
			manager.Broadcast(payload.GameID, mustMarshal("startGame", game))

		case "makeMove":
			var payload MakeMovePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}
			game, _ := GetGameByID(payload.GameID)
			if game == nil {
				fmt.Println("Game not found")
				continue
			}
			if err := utils.MakeMove(game, payload.PlayerID, payload.Cards); err != nil {
				fmt.Println("Invalid move:", err)
				//TODO: we should notify the user back also r
			}			
			manager.Broadcast(payload.GameID, mustMarshal("moveMade", game))

		case "check":
			var payload CheckPayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}
			result := handleCheck(payload.GameID, payload.CheckerId)
			manager.Broadcast(payload.GameID, mustMarshal("checkResult", result))

		case "burn":
			var payload StartGamePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}
			game := handleBurn(payload.GameID)
			manager.Broadcast(payload.GameID, mustMarshal("burned", game))

		default:
			fmt.Println("Unknown message type:", wsMsg.Type)
			respondText(conn, "unknown")
		}
	}
}


func respondText(conn *websocket.Conn, msg string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		fmt.Println("Write error:", err)
	}
}

func respondJSON(conn *websocket.Conn, msgType string, payload any) {
	resp := map[string]any{"type": msgType, "payload": payload}
	data, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, data)
}

func mustMarshal(msgType string, payload any) []byte {
	resp := map[string]any{"type": msgType, "payload": payload}
	data, _ := json.Marshal(resp)
	return data
}

func main() {
	fmt.Println("WebSocket server started at :8080")
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		handleWS(w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
