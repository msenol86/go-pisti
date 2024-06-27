package main

import (
	"fmt"
	"os"
	"time"
)

func fromGame(g Game, isPlayer1 bool, openCardCount int, isPlayerTurn bool, isGameOver bool) NetworkMessage {
	var player Player
	var opponent Player
	if isPlayer1 {
		player = g.player1
		opponent = g.player2
	} else {
		player = g.player2
		opponent = g.player1
	}

	gom := GameOverMessage{IsGameOver: false, PlayerWonCards: []Card{}, OpponentWonCards: []Card{}}
	if isGameOver {
		gom = GameOverMessage{IsGameOver: true, PlayerWonCards: player.wonCards, OpponentWonCards: opponent.wonCards}
	}

	return NetworkMessage{g.getBoardOpenCards(openCardCount), uint8(len(g.board)), uint8(len(g.deck)),
		player.hand, uint8(len(player.wonCards)), uint8(player.pistiCount), uint8(player.points),
		uint8(len(opponent.hand)), uint8(len(opponent.wonCards)), uint8(opponent.pistiCount), uint8(opponent.points), isPlayerTurn, gom}
}

type NetworkState struct {
	player1_addr string
	player2_addr string
}

func main() {
	tType := os.Getenv("TYPE")
	if tType == CLIENT {
		runClientGui()
	} else {
		player1Channel := make(chan NetworkMessage)
		player2Channel := make(chan NetworkMessage)
		player1InputChannel := make(chan int)
		player2InputChannel := make(chan int)
		playerJoinChannel := make(chan string)
		go startServer(player1Channel, player2Channel, playerJoinChannel, player1InputChannel, player2InputChannel)
		ns := NetworkState{"", ""}
		fmt.Println("Waiting for players to join")
		for {
			msg := <-playerJoinChannel
			if ns.player1_addr == "" {
				fmt.Printf("Adding %s as player 1\n", msg)
				ns.player1_addr = msg
			} else if ns.player2_addr == "" {
				fmt.Printf("Adding %s as player 2\n", msg)
				ns.player2_addr = msg
			} else {
				fmt.Println("All player slots are full!")
			}

			if ns.player1_addr != "" && ns.player2_addr != "" {
				break
			} else {
				time.Sleep(2 * time.Second)
			}
		}
		fmt.Println("Players joined")
		fmt.Printf("Player 1 %s\n", ns.player1_addr)
		fmt.Printf("Player 2 %s\n", ns.player2_addr)
		g := createAndStartGame()
		isGameOver := false
		isFirstRound := true
		for !isGameOver {
			var openCardCount int = 3

			if isFirstRound {
				openCardCount = 1
				player1Channel <- fromGame(g, true, openCardCount, false, isGameOver)
				player2Channel <- fromGame(g, false, openCardCount, false, isGameOver)
				time.Sleep(1 * time.Second)
			}

			player1Channel <- fromGame(g, true, openCardCount, true, isGameOver)
			s1 := <-player1InputChannel
			for len(player1InputChannel) > 0 {
				<-player1InputChannel
			}

			if isFirstRound {
				openCardCount = 2
			}

			g = g.playCard(true, s1-1)
			player1Channel <- fromGame(g, true, openCardCount, false, isGameOver)
			player2Channel <- fromGame(g, false, openCardCount, true, isGameOver)
			s2 := <-player2InputChannel
			for len(player2InputChannel) > 0 {
				<-player2InputChannel
			}

			if isFirstRound {
				openCardCount = 3
			}
			isFirstRound = false

			g = g.playCard(false, s2-1)
			player2Channel <- fromGame(g, false, openCardCount, false, isGameOver)
			if len(g.deck) > 7 && len(g.player1.hand) < 1 && len(g.player2.hand) < 1 {
				g = g.handOverCards(false)
				player1Channel <- fromGame(g, true, openCardCount, false, isGameOver)
				player2Channel <- fromGame(g, false, openCardCount, false, isGameOver)
				time.Sleep(1 * time.Second)
			}
			if g.lastCardPlayed() {
				// isGameOver = true
				// isFirstRound = true
				player1Channel <- fromGame(g, true, openCardCount, false, isGameOver)
				player2Channel <- fromGame(g, false, openCardCount, false, isGameOver)
				time.Sleep(1 * time.Second)
				isGameOver = true
				player1Channel <- fromGame(g, true, openCardCount, false, isGameOver)
				player2Channel <- fromGame(g, false, openCardCount, true, isGameOver)
			}
		}
		for {
			msg := <-playerJoinChannel
			fmt.Println(msg)
		}
	}
}
