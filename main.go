package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coderSomya/bluff/models"
	"github.com/coderSomya/bluff/utils"
	"github.com/gorilla/websocket"
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
	defer func() {
		manager.RemoveClient(conn) // Clean up on disconnect
		conn.Close()
	}()

	var gameID string
	var playerID string // Track current player

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

			manager.AddClient(gameID, conn, playerID)

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
			playerID = payload.Player.PlayerId

			manager.AddClient(gameID, conn, playerID)
			
			// Send personalized game state to all players
			updatedGame, _ := GetGameByID(gameID)
			broadcastPersonalizedGameState(gameID, updatedGame, "playerJoined", map[string]interface{}{
				"newPlayer": payload.Player.Name,
				"message": fmt.Sprintf("%s joined the game", payload.Player.Name),
			})

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
				respondJSON(conn, "error", map[string]string{"message": "Game not found"})
				continue
			}
			
			if err := utils.MakeMove(game, payload.PlayerID, payload.Cards); err != nil {
				respondJSON(conn, "error", map[string]string{"message": err.Error()})
				continue
			}
			
			// Broadcast personalized game state
			eventData := map[string]interface{}{
				"action": "move",
				"playerId": payload.PlayerID,
				"cardCount": len(payload.Cards),
				"claimedCard": game.MoveCard,
				"message": fmt.Sprintf("Player %s played %d cards", payload.PlayerID, len(payload.Cards)),
			}
			
			broadcastPersonalizedGameState(payload.GameID, game, "moveMade", eventData)

		case "check":
			var payload CheckPayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}
			
			game, _ := GetGameByID(payload.GameID)
			if game == nil {
				respondJSON(conn, "error", map[string]string{"message": "Game not found"})
				continue
			}
			
			// Capture state before check
			lastPlayer := ""
			if game.LastPlayerId != nil {
				lastPlayer = *game.LastPlayerId
			}
			cardsInPile := len(game.PlayedCards)
			
			// Perform the check
			wasHonest, err := utils.Check(game, payload.CheckerId)
			if err != nil {
				respondJSON(conn, "error", map[string]string{"message": err.Error()})
				continue
			}
			
			// Create event data
			eventData := map[string]interface{}{
				"checker": payload.CheckerId,
				"checkedPlayer": lastPlayer,
				"wasBluff": !wasHonest,
				"cardsTransferred": cardsInPile,
			}
			
			if wasHonest {
				eventData["result"] = fmt.Sprintf("%s was honest! %s takes the cards", lastPlayer, payload.CheckerId)
			} else {
				eventData["result"] = fmt.Sprintf("%s was bluffing! %s takes the cards", lastPlayer, lastPlayer)
			}
			
			// Broadcast personalized updates
			broadcastPersonalizedGameState(payload.GameID, game, "checkResult", eventData)

		case "burn":
			var payload StartGamePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}
			
			game, _ := GetGameByID(payload.GameID)
			if game == nil {
				respondJSON(conn, "error", map[string]string{"message": "Game not found"})
				continue
			}
			
			// Validate burn is possible
			if !utils.IsBurnPossible(*game) {
				respondJSON(conn, "error", map[string]string{
					"message": "Burn not possible - you must be the last player who made a move",
				})
				continue
			}
			
			// Get current player (who's trying to burn)
			burnerID := manager.GetPlayerID(conn)
			if burnerID == "" {
				respondJSON(conn, "error", map[string]string{"message": "Player not found"})
				continue
			}
			
			// Capture cards before burn
			cardsBurned := len(game.PlayedCards)
			
			// Perform burn
			*game = handleBurn(payload.GameID)
			
			// Create event data
			eventData := map[string]interface{}{
				"burner": burnerID,
				"cardsBurned": cardsBurned,
				"message": fmt.Sprintf("Player %s burned the pile!", burnerID),
			}
			
			// Broadcast personalized updates
			broadcastPersonalizedGameState(payload.GameID, game, "burned", eventData)

		default:
			fmt.Println("Unknown message type:", wsMsg.Type)
			respondText(conn, "unknown")
		}
	}
}

// Create a safe game state for a specific player
func createPlayerGameState(game *models.Game, forPlayerID string) map[string]interface{} {
	playerCards := []models.Card{}
	
	// Find this player's cards
	for _, player := range game.Players {
		if player.PlayerId == forPlayerID {
			playerCards = player.Cards
			break
		}
	}
	
	// Create safe player list (without cards)
	safePlayers := make([]map[string]interface{}, len(game.Players))
	for i, player := range game.Players {
		safePlayers[i] = map[string]interface{}{
			"playerId": player.PlayerId,
			"name": player.Name,
			"cardCount": len(player.Cards),
		}
	}
	
	return map[string]interface{}{
		"gameId": game.GameId,
		"players": safePlayers,
		"yourCards": playerCards,
		"currentPlayerId": game.CurrentPlayerId,
		"playedCardsCount": len(game.PlayedCards),
		"lastPlayerId": game.LastPlayerId,
		"lastPlayedQty": game.LastPlayedQty,
		"moveCard": game.MoveCard,
	}
}

// Broadcast personalized messages to all players in a game
func broadcastPersonalizedGameState(gameID string, game *models.Game, eventType string, eventData map[string]interface{}) {
	manager.BroadcastPersonalized(gameID, func(playerID string) []byte {
		safeGameState := createPlayerGameState(game, playerID)
		
		response := map[string]interface{}{
			"type": eventType,
			"payload": map[string]interface{}{
				"game": safeGameState,
				"event": eventData,
			},
		}
		
		data, _ := json.Marshal(response)
		return data
	})
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


