package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins (for development)
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

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
			fmt.Println("New Game:", payload.Creator)
			// Respond with confirmation
			respondText(conn, wsMsg.Type)

		case "newPlayer":
			var payload NewPlayerPayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				fmt.Println("Invalid newPlayer payload:", err)
				continue
			}
			fmt.Println("New Player:", payload.Player, "to Game ID:", payload.GameID)
			respondText(conn, wsMsg.Type)

		case "startGame":
			var payload StartGamePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				fmt.Println("Invalid startGame payload:", err)
				continue
			}
			fmt.Println("Start Game:", payload.GameID)
			respondText(conn, wsMsg.Type)

		case "makeMove":
			var payload MakeMovePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				fmt.Println("Invalid makeMove payload:", err)
				continue
			}
			fmt.Printf("Player %s played %+v in Game %s\n", payload.PlayerID, payload.Cards, payload.GameID)
			respondText(conn, wsMsg.Type)

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

func main() {
	fmt.Println("WebSocket server started at :8080")
	http.HandleFunc("/ws", handleWS)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
