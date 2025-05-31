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

func Check(game *models.Game, checkerId string) (bool, error) {
	if game.LastPlayerId == nil || game.LastPlayedQty == nil {
		return false, fmt.Errorf("no previous move to check")
	}

	lastPlayerId := *game.LastPlayerId
	qty := *game.LastPlayedQty

	if qty > len(game.PlayedCards) {
		return false, fmt.Errorf("invalid state: more cards played than exist")
	}

	startIdx := len(game.PlayedCards) - qty
	lastPlayed := game.PlayedCards[startIdx:]

	moveCard := game.MoveCard
	if moveCard == nil {
		return false, fmt.Errorf("no move card to check against")
	}

	bluffed := false
	for _, card := range lastPlayed {
		if card.Value != moveCard.Value && card.Value != models.Joker {
			bluffed = true
			break
		}
	}

	receiverId := checkerId
	if bluffed {
		receiverId = lastPlayerId
	}

	for i := range game.Players {
		if game.Players[i].PlayerId == receiverId {
			game.Players[i].Cards = append(game.Players[i].Cards, game.PlayedCards...)
			break
		}
	}

	game.PlayedCards = []models.Card{}
	game.LastPlayerId = nil
	game.LastPlayedQty = nil
	game.MoveCard = nil

	return !bluffed, nil // true = honest move, false = bluff
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