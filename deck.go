package main

import (
	"fmt"
	"math/rand"
)

type Deck []Card

func (d Deck) print() {
	for i, card := range d {
		fmt.Println(i, card)
	}
}

func createDeck() Deck {
	myDeck := []Card{}
	for _, tSuit := range SUITS {
		for _, tRank := range RANKS {
			myDeck = append(myDeck, Card{SuitType(tSuit), RankType(tRank)})
		}
	}
	return myDeck
}

func (d Deck) shuffleDeck() {
	// rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}
