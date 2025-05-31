package main

import (
	"github.com/coderSomya/bluff/models"
	"encoding/json"
)

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // defer decoding until we know the type
}

type NewGamePayload struct {
	Creator models.Player `json:"creator"`
}

type NewPlayerPayload struct {
	GameID string        `json:"gameId"`
	Player models.Player `json:"player"`
}

type StartGamePayload struct {
	GameID string `json:"gameId"`
}

type MakeMovePayload struct {
	GameID   string        `json:"gameId"`
	PlayerID string        `json:"playerId"`
	Cards    []models.Card `json:"cards"`
}


type WSResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"error,omitempty"`
}
