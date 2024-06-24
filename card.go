package main

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
	suit SuitType
	rank RankType
}

func RemoveIndex(s []Card, index int) []Card {
	return append(s[:index], s[index+1:]...)
}
