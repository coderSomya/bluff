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
		// Allow all origins (for development)
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var gameID string

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
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
			addPlayerToGame(payload.GameID, payload.Player)
			gameID = payload.GameID
			manager.AddClient(gameID, conn)
			manager.Broadcast(gameID, []byte(fmt.Sprintf("Player %s joined", payload.Player.PlayerId)))

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
			var payload StartGamePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}
			result := handleCheck(payload.GameID)
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
	http.HandleFunc("/ws", handleWS)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
