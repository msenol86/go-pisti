package main

import (
	"fmt"
)

type Player struct {
	hand       []Card
	wonCards   []Card
	pistiCount int
	points     int
}

type Game struct {
	deck    Deck
	board   []Card
	player1 Player
	player2 Player // ai player in single player games
}

func createAndStartGame() Game {
	player1 := Player{[]Card{}, []Card{}, 0, 0}
	player2 := Player{[]Card{}, []Card{}, 0, 0}
	g := Game{createDeck(), []Card{}, player1, player2}
	g.deck.shuffleDeck()
	g = g.handOverCards(true)
	return g
}

func (g Game) handOverCards(isNewGame bool) Game {
	for i := 0; i < 4; i++ {
		if isNewGame {
			g.board = append(g.board, g.deck[0])
			g.deck = g.deck[1:]
		}
		g.player1.hand = append(g.player1.hand, g.deck[0])
		g.deck = g.deck[1:]
		g.player2.hand = append(g.player2.hand, g.deck[0])
		g.deck = g.deck[1:]
	}
	return g
}

func (g Game) playCard(isPlayer1Card bool, cardIndex int) Game {
	var selectedCard Card
	if isPlayer1Card {
		selectedCard = g.player1.hand[cardIndex]
		g.player1.hand = RemoveIndex(g.player1.hand, cardIndex)
	} else {
		selectedCard = g.player2.hand[cardIndex]
		g.player2.hand = RemoveIndex(g.player2.hand, cardIndex)
	}
	if isPlayer1Card {
		fmt.Println("player1 played card:", selectedCard)
	} else {
		fmt.Println("player2 played card:", selectedCard)
	}

	g.board = append(g.board, selectedCard)
	if len(g.board) > 1 {
		if g.lastCardPlayed() || g.board[len(g.board)-1].Rank == J || g.board[len(g.board)-1].Rank == g.board[len(g.board)-2].Rank {
			fmt.Println("Win!\a")
			for i := 0; i < len(g.board); i++ {
				if isPlayer1Card {
					g.player1.wonCards = append(g.player1.wonCards, g.board[i])
				} else {
					g.player2.wonCards = append(g.player2.wonCards, g.board[i])
				}

			}
			g.board = []Card{}
		}
	}
	return g
}

func (g Game) boardToStringSlice() []string {
	db := []string{}
	for i := 0; i < len(g.board); i++ {
		var astr string
		if i == len(g.board)-1 {
			astr = string(g.board[i].Suit) + fmt.Sprint(g.board[i].Rank)
		} else {
			astr = "*-*"
		}

		db = append(db, astr)
	}
	return db
}

func (g Game) getBoardOpenCards(openCardCount int) []Card {
	cards := []Card{}
	for i := openCardCount; i > 0; i-- {
		if len(g.board)-i > -1 {
			cards = append(cards, g.board[len(g.board)-i])
		}
	}
	return cards
}

func (g Game) lastCardPlayed() bool {
	if len(g.deck) == 0 && len(g.player1.hand) == 0 && len(g.player2.hand) == 0 {
		return true
	} else {
		return false
	}
}
