
package utils

import (
	"github.com/coderSomya/bluff/models"
	"fmt"
	"math/rand"
	"time"
)


// create the initial deck
func InitializeDeck() []models.Card {

	var deck = make([]models.Card, 0)

	var AllSuits = []models.Suit{models.Diamonds, models.Spades, models.Clubs, models.Hearts}
	var AllRanks = []models.Rank{models.Ace, models.Two, models.Three,
								 models.Four, models.Five, models.Six,
								 models.Seven, models.Eight, models.Nine, 
								 models.Ten, models.Jack, models.Queen, models.King}

	for _, suit := range AllSuits {
		for _, rank := range AllRanks {
			deck = append(deck, models.Card{Value: rank, Type: suit})
		}
	}

	deck = append(deck, models.Card{Value: models.Joker, Type: models.Undefined}, models.Card{Value: models.Joker,Type: models.Undefined}, models.Card{Value: models.Joker, Type: models.Undefined})

	for _, card := range deck {
		fmt.Println(card)
	}

	return deck
}

func RandomizeDeck(n int) [][]models.Card {
		if n <= 0 {
		return nil
	}

	deck := InitializeDeck() // pura set le aao

	// generate a random number
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	groupSize := len(deck) / n

	result := make([][]models.Card, 0, n)

	for i := 0; i < n; i++ {
		start := i * groupSize
		end := start + groupSize
		group := deck[start:end]
		result = append(result, group)
	}

	// i am going to remove the remainder cards

	for i, group := range result {
		fmt.Printf("Group %d (%d cards):\n", i+1, len(group))
		for _, card := range group {
			fmt.Println(card)
		}
		fmt.Println()
	}

	return result
}


