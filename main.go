package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func startServer(data chan int) {
	// Bind the port.
	ServerAddr, err := net.ResolveUDPAddr("udp", ":6666")
	if err != nil {
		fmt.Println("Error binding port!")
	}

	ServerConn, _ := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	for {
		// Recieve a UDP packet and unmarshal it into a protobuf.
		n, adrr, _ := ServerConn.ReadFromUDP(buf)
		fmt.Println("Packet received!", n, adrr)
		// data <- n
		// print(data)
		// Do stuff with buf.
	}
}

func main() {
	channel := make(chan int)
	go startServer(channel)
	g := createAndStartGame()
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
