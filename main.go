package main

import (
	"fmt"
	"strconv"
	"time"
)

type NetworkState struct {
	player1_addr string
	player2_addr string
}

func main() {
	channel := make(chan string)
	playerJoinChannel := make(chan string)
	go startServer(channel, playerJoinChannel)
	g := createAndStartGame()
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
	isGameOver := false
	for !isGameOver {
		fmt.Println("playerHand:", g.player1.hand)
		fmt.Println("board", g.boardToStringSlice())
		fmt.Printf("Input the number of the card do you want to play? [1-%d]: ", len(g.player1.hand))
		var input string
		fmt.Scanln(&input)
		if s, err := strconv.Atoi(input); err == nil {
			// fmt.Printf("%T, %v", s, s-1)
			g = g.playCard(true, s-1)
		}
		fmt.Print("\nAI playing card")
		for i := 0; i < 2; i++ {
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
		g = g.playCard(false, 0)
		fmt.Println("--------------\n\n")
		if len(g.deck) > 7 && len(g.player1.hand) < 1 && len(g.player2.hand) < 1 {
			fmt.Print("Handin over new cards")
			for i := 0; i < 3; i++ {
				time.Sleep(1 * time.Second)
				fmt.Print(".")
			}
			g = g.handOverCards(false)
		}
		if len(g.deck) == 0 && len(g.player1.hand) == 0 && len(g.player2.hand) == 0 {
			isGameOver = true
		}
	}
}
