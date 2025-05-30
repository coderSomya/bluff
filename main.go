package main

import "fmt"

import (
	"github.com/coderSomya/bluff/utils"
)

func main() {
    fmt.Println("Hello, Go!")

	utils.InitializeDeck()
	
	utils.RandomizeDeck(4)
}
