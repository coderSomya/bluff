package main

import (
	"github.com/coderSomya/bluff/models"
	"github.com/coderSomya/bluff/utils"
)

var gameManager := models.GameManager{
	games: make([]models.Game,0)
}


/*
*/

func addNewGame(creator Player){
	
	var newGame := models.Game{
		GameId: "123"
		Players: make([]models.Player,0)
		CurrentPlayerId: creator.PlayerId
		PlayedCards: make([]models.Card,0)
		LastPlayerId: nil
		LastPlayedQty: nil
		MoveCard: nil
	}
	newGame.Players = append(newGame.Players, creator)
	gameManager = gameManager.append(newGame)
}

func addPlayerToGame(gameId string, player Player){
	for i, game := range gameManager.games {
		if game.GameId == gameId {
			gameManager.games[i].Players = append(gameManager[i].Players, player)
			return
		}
	}
	fmt.Println("Game not found:", gameId)
}

