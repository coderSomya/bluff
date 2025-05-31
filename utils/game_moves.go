package utils

import (
	"github.com/coderSomya/bluff/models"
)

func IsGameOver(game models.Game) bool{

	if len(game.PlayedCards)!=0{
		return false
	}

	for i := range len(game.Players){
		if len(game.Players[i].Cards) == 0 {
			return true
		}
	}
	return false
}

func IsBurnPossible(game models.Game) bool{

	if game.LastPlayerId == game.CurrentPlayerId {
		return true
	}
	
	return false
}

func Check(game models.Game) bool{
	if game.CurrentPlayerId == game.LastPlayerId {
		return false
	}

	if len(game.PlayedCards) == 0 {
		return false
	}

	var qty = game.LastPlayedQty //itne cards dekhne hain

	for i := len(game.PlayedCards) - 1; i >= *qty; i-- {
		if game.PlayedCards[i]!= *game.MoveCard && game.PlayedCards[i].Value != models.Joker {
			return true
		}
	}

	return false
}

func MakeMove(game models.Game, playerId string, cards []models.Card) {
	
	var qty = len(cards)

	for i := 0; i<qty; i++{
		game.PlayedCards = append(game.PlayedCards, cards[i])
	}

	game.LastPlayedQty = &qty
	game.LastPlayerId = &playerId
}