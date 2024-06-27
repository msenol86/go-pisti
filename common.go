package main

import (
	"fmt"
)

const (
	HOST   = "localhost"
	PORT   = "6666"
	TYPE   = "tcp"
	JOIN   = "JOIN"
	PLAY   = "PLAY"
	CLIENT = "CLIENT"
	AI     = "AI"
)
const UNKNOWN = "� �"

type SuitType string
type RankType int

var SUITS = [4]string{"♠", "♥", "♦", "♣"}
var RANKS = [13]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

const (
	SPADES   SuitType = "♠"
	HEARTS   SuitType = "♥"
	DIAMONDS SuitType = "♦"
	CLUBS    SuitType = "♣"
)

const (
	ACE   RankType = 1
	TWO   RankType = 2
	THREE RankType = 3
	FOUR  RankType = 4
	FIVE  RankType = 5
	SIX   RankType = 6
	SEVEN RankType = 7
	EIGHT RankType = 8
	NINE  RankType = 9
	TEN   RankType = 10
	J     RankType = 11
	Q     RankType = 12
	K     RankType = 13
)

type Card struct {
	Suit SuitType
	Rank RankType
}

func (r RankType) toString() string {
	switch r {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		return fmt.Sprint(r)
	}
}

func (c Card) toString() string {

	suit := fmt.Sprint(c.Suit)
	return suit + " " + c.Rank.toString()
}

func RemoveIndex(s []Card, index int) []Card {
	return append(s[:index], s[index+1:]...)
}

type GameOverMessage struct {
	IsGameOver       bool
	PlayerWonCards   []Card
	OpponentWonCards []Card
}

type NetworkMessage struct {
	BoardOpenCards        []Card
	BoardCount            uint8
	DeckCount             uint8
	PlayerHand            []Card
	PlayerWonCardsCount   uint8
	PlayerPistiCounts     uint8
	PlayerPoints          uint8
	OpponentHandCount     uint8
	OpponentWonCardsCount uint8
	OpponentPistiCounts   uint8
	OpponentPoints        uint8
	IsPlayerTurn          bool
	GameOverState         GameOverMessage
}
