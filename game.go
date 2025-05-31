package main

import (
	"github.com/coderSomya/bluff/models"
	"github.com/coderSomya/bluff/utils"
	"fmt"
)

var gameManager = models.GameManager{
	Games: make([]models.Game,0),
}


/*
*/

func GetGameByID(gameID string) (*models.Game, int) {
	for i, game := range gameManager.Games {
		if game.GameId == gameID {
			return &gameManager.Games[i], i
		}
	}
	return nil, -1
}

func addNewGame(creator models.Player) models.Game{
	var newGame = models.Game{
		GameId: "123",
		Players: make([]models.Player,0),
		CurrentPlayerId: nil,
		PlayedCards: make([]models.Card,0),
		LastPlayerId: nil,
		LastPlayedQty: nil,
		MoveCard: nil,
	}
	newGame.Players = append(newGame.Players, creator)
	newGame.CurrentPlayerId = &creator.PlayerId
	gameManager.AddGame(newGame)

	return newGame
}

func addPlayerToGame(gameId string, player models.Player) (*models.Game, error){
	for i, game := range gameManager.Games {
		if game.GameId == gameId {
			gameManager.Games[i].Players = append(gameManager.Games[i].Players, player)
			return &gameManager.Games[i], nil
		}
	}
	return nil, fmt.Errorf("game not found: %s", gameId)
}

func handleCheck(gameID string) map[string]any {
	for _, game := range gameManager.Games {
		if game.GameId == gameID {
			result := utils.Check(game)
			return map[string]any{
				"result": result,
			}
		}
	}
	return map[string]any{"error": "game not found"}
}

func handleBurn(gameID string) models.Game {
	for i, game := range gameManager.Games {
		if game.GameId == gameID && utils.IsBurnPossible(game) {
			game.PlayedCards = []models.Card{}
			gameManager.Games[i] = game
			return game
		}
	}
	return models.Game{}
}

func createGame(creator models.Player, gameID string) models.Game {
	newGame := models.Game{
		GameId:         gameID,
		Players:        []models.Player{creator},
		CurrentPlayerId: &creator.PlayerId,
		PlayedCards:    []models.Card{},
	}
	gameManager.AddGame(newGame)
	return newGame
}

func startGame(gameID string) models.Game {
	for i, game := range gameManager.Games {
		if game.GameId == gameID {
			return gameManager.Games[i]
		}
	}
	return models.Game{}
}
