package main

import (
	"fmt"
	"strconv"
)

func main() {
	// theDeck := createDeck()
	// theDeck.shuffleDeck()
	// theDeck.print()

	// board := make([]Card, 0)
	// playerHand := make([]Card, 0)
	// opponentHand := make([]Card, 0)
	// playerWonCards := make([]Card, 0)
	// opponentWonCards := make([]Card, 0)
	// playerPistiCount := 0
	// opponentPistiCount := 0
	// playerPoints := 0
	// opponentPoints := 0
	g := createAndStartGame()
	g.deck.print()
	fmt.Println("playerHand:", g.playerHand)
	fmt.Println("board", g.boardToStringSlice())

	fmt.Print("Input the number of the card do you want to play? [1-4]: ")
	var input string
	fmt.Scanln(&input)
	if s, err := strconv.Atoi(input); err == nil {
		fmt.Printf("%T, %v", s, s-1)
		g = g.playCard(true, s-1)
		fmt.Println("playerHand:", g.playerHand)
		fmt.Println("board", g.boardToStringSlice())
	}
	// fmt.Print(input)

	// g = g.playCard(true, age)

	// fmt.Println("playerHand:", g.playerHand)
	// fmt.Println("board", g.boardToStringSlice())

	// xx := cardIndex - 1
	// fmt.Println("xx", xx)
	// g = g.playCard(true, xx)

	// g = g.playCard(true, 0)
	// fmt.Println("playerHand:", g.playerHand)
	// fmt.Println("board:", g.boardToStringSlice())
	// g = g.playCard(false, 0)
	// fmt.Println("playerHand:", g.playerHand)
	// fmt.Println("board:", g.boardToStringSlice())
	// g = g.playCard(true, 0)
	// fmt.Println("playerHand:", g.playerHand)
	// fmt.Println("board:", g.boardToStringSlice())
	// g = g.playCard(false, 0)
	// fmt.Println("playerHand:", g.playerHand)
	// fmt.Println("board:", g.boardToStringSlice())
	// g = g.playCard(true, 0)
	// fmt.Println("playerHand:", g.playerHand)
	// fmt.Println("board:", g.boardToStringSlice())
}

// func startGame(pDeck Deck, pBoard []Card, pPlayerHand []Card, popponentHand []Card, pPlayerPistiCount int, pOpponentPistiCount int, pPlayerPoints int, pOpponentPoints int) {
// 	for i := 0; i < 5; i++ {
// 		pBoard = append(pBoard, pDeck[0])
// 		pDeck = pDeck[1:]
// 		pPlayerHand = append(pPlayerHand, pDeck[0])
// 		pDeck = pDeck[1:]
// 		popponentHand = append(popponentHand, pDeck[0])
// 		pDeck = pDeck[1:]
// 	}
// 	pPlayerPistiCount = 0
// }
