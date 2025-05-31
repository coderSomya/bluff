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

func handleCheck(gameID string, checkerId string) map[string]any {
	for _, game := range gameManager.Games {
		if game.GameId == gameID {
			ok, err := utils.Check(&game, checkerId)
			if err != nil {
				// TODO: return error to client
			}
			if ok {
				// TODO: inform that move was honest, checker takes cards
			} else {
				// TODO: inform that move was a bluff, last player takes cards
			}
			return map[string]any{
				"result": ok,
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
			numPlayers := len(game.Players)
			if numPlayers == 0 {
				fmt.Println("No players in game")
				return game
			}

			// Randomly distribute kardo deck
			hands := utils.RandomizeDeck(numPlayers)

			for j := 0; j < numPlayers; j++ {
				gameManager.Games[i].Players[j].Cards = hands[j]
			}

			// basically creator will be first guy to make a move
			firstPlayerID := gameManager.Games[i].Players[0].PlayerId
			gameManager.Games[i].CurrentPlayerId = &firstPlayerID

			fmt.Printf("Game %s started with %d players\n", gameID, numPlayers)
			return gameManager.Games[i]
		}
	}
	return models.Game{}
}

