package main

import (
	"github.com/coderSomya/bluff/models"
	"fmt"
)

var gameManager = models.GameManager{
	Games: make([]models.Game,0),
}


/*
*/

func addNewGame(creator models.Player){
	var newGame = models.Game{
		GameId: "123",
		Players: make([]models.Player,0),
		CurrentPlayerId: &creator.PlayerId,
		PlayedCards: make([]models.Card,0),
		LastPlayerId: nil,
		LastPlayedQty: nil,
		MoveCard: nil,
	}
	newGame.Players = append(newGame.Players, creator)
	gameManager.AddGame(newGame)
}

func addPlayerToGame(gameId string, player models.Player){
	for i, game := range gameManager.Games {
		if game.GameId == gameId {
			gameManager.Games[i].Players = append(gameManager.Games[i].Players, player)
			return
		}
	}
	fmt.Println("Game not found:", gameId)
}

