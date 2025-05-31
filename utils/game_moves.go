package utils

import (
	"github.com/coderSomya/bluff/models"
	"fmt"
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

func MakeMove(game *models.Game, playerId string, cards []models.Card) error {
	playerIndex := -1
	for i, player := range game.Players {
		if player.PlayerId == playerId {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return fmt.Errorf("player with id %s not found in game", playerId)
	}

	playerCards := game.Players[playerIndex].Cards
	playedCardsMap := make(map[string]int)

	for _, card := range cards {
		key := card.Type.String() + "-" + card.Value.String()
		playedCardsMap[key]++
	}

	playerCardCount := make(map[string]int)
	for _, card := range playerCards {
		key := card.Type.String() + "-" + card.Value.String()
		playerCardCount[key]++
	}

	for key, count := range playedCardsMap {
		if playerCardCount[key] < count {
			return fmt.Errorf("player does not own enough of card %s", key)
		}
	}

	remaining := []models.Card{}
	cardNeeded := make(map[string]int)
	for _, card := range cards {
		key := card.Type.String() + "-" + card.Value.String()
		cardNeeded[key]++
	}

	for _, card := range playerCards {
		key := card.Type.String() + "-" + card.Value.String()
		if cardNeeded[key] > 0 {
			cardNeeded[key]--
		} else {
			remaining = append(remaining, card)
		}
	}

	game.Players[playerIndex].Cards = remaining

	for _, card := range cards {
		game.PlayedCards = append(game.PlayedCards, card)
	}
	qty := len(cards)
	game.LastPlayedQty = &qty
	game.LastPlayerId = &playerId

	return nil
}

/*
type Game struct {
    GameId  string
    Players []Player
	CurrentPlayerId *string
	PlayedCards []Card
	LastPlayerId *string
	LastPlayedQty *int
	MoveCard *Card
}
*/

func Pass(game *models.Game){
	var currentPlayerId = game.CurrentPlayerId;
	var numPlayers = len(game.Players);
	var last_player_index = -1;
	for i:=0; i<numPlayers; i++{
		if game.Players[i].PlayerId == *currentPlayerId {
			last_player_index = i
		}
	}
	
	game.CurrentPlayerId = &game.Players[(last_player_index+1)%numPlayers].PlayerId;
}