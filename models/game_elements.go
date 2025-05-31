package models

type GameManager struct{
	games []Game
}

type Game struct {
    GameId  string
    Players []Player
	CurrentPlayerId *string
	PlayedCards []Card
	LastPlayerId *int
	LastPlayedQty *int
	MoveCard *Card
}

type Player struct {
	PlayerId string
	Name string
	Cards []Card
}


type Card struct {
	Value Rank
	Type Suit
}

type Suit int

const (
    Diamonds Suit = iota
    Spades
    Clubs
    Hearts
	Undefined
)

type Rank int

const (
    Ace Rank = iota+1
    Two
    Three
    Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Joker
)


func (s Suit) String() string {
	return [...]string{"Spades", "Hearts", "Diamonds", "Clubs", "Undefined"}[s]
}

func (r Rank) String() string {
	return [...]string{
		"Ace", "Two", "Three", "Four", "Five", "Six", "Seven",
		"Eight", "Nine", "Ten", "Jack", "Queen", "King", "Joker",
	}[r-1]
}

func (c Card) String() string {
	return c.Value.String() + " of " + c.Type.String()
}
