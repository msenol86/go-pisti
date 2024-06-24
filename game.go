package main

import (
	"fmt"
)

type Game struct {
	deck               Deck
	board              []Card
	playerHand         []Card
	opponentHand       []Card
	playerWonCards     []Card
	opponentWonCards   []Card
	playerPistiCount   int
	opponentPistiCount int
	playerPoints       int
	opponentPoints     int
}

func createAndStartGame() Game {
	g := Game{createDeck(), []Card{}, []Card{}, []Card{}, []Card{}, []Card{}, 0, 0, 0, 0}
	g.deck.shuffleDeck()
	for i := 0; i < 4; i++ {
		g.board = append(g.board, g.deck[0])
		g.deck = g.deck[1:]
		g.playerHand = append(g.playerHand, g.deck[0])
		g.deck = g.deck[1:]
		g.opponentHand = append(g.opponentHand, g.deck[0])
		g.deck = g.deck[1:]
	}
	return g
}

func (g Game) playCard(isPlayerCard bool, cardIndex int) Game {
	var selectedCard Card
	if isPlayerCard {
		selectedCard = g.playerHand[cardIndex]
		g.playerHand = RemoveIndex(g.playerHand, cardIndex)
	} else {
		selectedCard = g.opponentHand[cardIndex]
		g.opponentHand = RemoveIndex(g.opponentHand, cardIndex)
	}
	if isPlayerCard {
		fmt.Println("player played card:", selectedCard)
	} else {
		fmt.Println("opponent played card:", selectedCard)
	}

	g.board = append(g.board, selectedCard)
	return g
}

func (g Game) boardToStringSlice() []string {
	db := []string{}
	for i := 0; i < len(g.board); i++ {
		var astr string
		if i == len(g.board)-1 {
			astr = string(g.board[i].suit) + fmt.Sprint(g.board[i].rank)
		} else {
			astr = "*-*"
		}

		db = append(db, astr)
	}
	return db
	// if len(g.board) == 0 {
	// 	fmt.Println("board: []")
	// }
	// for i := 0; i < len(g.board); i++ {
	// 	if i == len(g.board)-1 {

	// 	} else {
	// 		fmt.Println(g.board[i])
	// 	}

	// }
}
