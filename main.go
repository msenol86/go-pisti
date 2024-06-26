package main

import (
	"fmt"
	"time"
)

// type Season int64

// const (
// 	Summer Season = iota
// 	Autumn
// 	Winter
// 	Spring
// )

// type Foo struct {
//     Number int    `json:"number"`
//     Title  string `json:"title"`
// }

// foo_marshalled, err := json.Marshal(Foo{Number: 1, Title: "test"})
// fmt.Fprint(w, string(foo_marshalled)) // write response to ResponseWriter (w)

type NetworkMessage struct {
	Board                 []string
	DeckCount             uint8
	PlayerHand            []Card
	PlayerWonCardsCount   uint8
	PlayerPistiCounts     uint8
	PlayerPoints          uint8
	OpponentHand          []Card
	OpponentWonCardsCount uint8
	OpponentPistiCounts   uint8
	OpponentPoints        uint8
}

func fromGame(g Game, isPlayer1 bool) NetworkMessage {
	var player Player
	var opponent Player
	if isPlayer1 {
		player = g.player1
		opponent = g.player2
	} else {
		player = g.player2
		opponent = g.player1
	}

	return NetworkMessage{g.boardToStringSlice(), uint8(len(g.deck)),
		player.hand, uint8(len(player.wonCards)), uint8(player.pistiCount), uint8(player.points),
		opponent.hand, uint8(len(opponent.wonCards)), uint8(opponent.pistiCount), uint8(opponent.points)}
}

type NetworkState struct {
	player1_addr string
	player2_addr string
}

func main() {
	player1Channel := make(chan NetworkMessage)
	player2Channel := make(chan NetworkMessage)
	playerJoinChannel := make(chan string)
	go startServer(player1Channel, player2Channel, playerJoinChannel)
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
	counter := 0
	for !isGameOver {
		// fmt.Println("playerHand:", g.player1.hand)
		// fmt.Println("board", g.boardToStringSlice())
		if counter < 1 {
			// player1Channel <- strings.Join(g.boardToStringSlice()[:], ",")
			fmt.Println("Game", g.player1.hand)
			player1Channel <- fromGame(g, true)
		}
		counter += 1

		// fmt.Printf("Input the number of the card do you want to play? [1-%d]: ", len(g.player1.hand))
		// var input string
		// fmt.Scanln(&input)
		// if s, err := strconv.Atoi(input); err == nil {
		// 	// fmt.Printf("%T, %v", s, s-1)
		// 	g = g.playCard(true, s-1)
		// }
		// fmt.Print("\nAI playing card")
		// for i := 0; i < 2; i++ {
		// 	time.Sleep(1 * time.Second)
		// 	fmt.Print(".")
		// }
		// g = g.playCard(false, 0)
		// fmt.Println("--------------\n\n")
		// if len(g.deck) > 7 && len(g.player1.hand) < 1 && len(g.player2.hand) < 1 {
		// 	fmt.Print("Handing over new cards")
		// 	for i := 0; i < 3; i++ {
		// 		time.Sleep(1 * time.Second)
		// 		fmt.Print(".")
		// 	}
		// 	g = g.handOverCards(false)
		// }
		// if len(g.deck) == 0 && len(g.player1.hand) == 0 && len(g.player2.hand) == 0 {
		// 	isGameOver = true
		// }
	}
}
