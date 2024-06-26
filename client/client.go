package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "6666"
	TYPE = "tcp"
)

type SuitType string
type RankType int

type Card struct {
	Suit SuitType
	Rank RankType
}

type NetworkMessage struct {
	Board                 []string
	DeckCount             uint8
	PlayerHand            []Card
	PlayerWonCardsCount   uint8
	PlayerPistiCounts     uint8
	PlayerPoints          uint8
	OpponentHandCount     uint8
	OpponentWonCardsCount uint8
	OpponentPistiCounts   uint8
	OpponentPoints        uint8
}

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	// _, err = conn.Write([]byte("This is a message\n"))
	_, err = conn.Write([]byte("JOIN\n"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}
	// st := `{"Board":["*-*","*-*","*-*","♣8"],"DeckCount":40,"PlayerHand":[{"Suit":"♠","Rank":9},{"Suit":"♦","Rank":8},{"Suit":"♠","Rank":8},{"Suit":"♦","Rank":2}],"PlayerWonCardsCount":0,"PlayerPistiCounts":0,"PlayerPoints":0,"OpponentHand":[{"Suit":"♦","Rank":3},{"Suit":"♣","Rank":9},{"Suit":"♦","Rank":7},{"Suit":"♦","Rank":6}],"OpponentWonCardsCount":0,"OpponentPistiCounts":0,"OpponentPoints":0}`
	// var nm NetworkMessage
	// if err := json.Unmarshal([]byte(st), &nm); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("nm", nm)

	for {
		// buffer to get data
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		// println("Received message:", string(received))
		// clean_msg := strings.TrimSpace(string(received))
		var nm NetworkMessage
		d := json.NewDecoder(conn)
		if err := d.Decode(&nm); err != nil {
			fmt.Println(err)
		}
		// fmt.Println("nm", nm)
		fmt.Println("Board: ", nm.Board)
		fmt.Println("Your Hand", nm.PlayerHand)
		fmt.Printf("Type the number of card you want to play [%d-%d]: ", 1, len(nm.PlayerHand))
		var input string
		fmt.Scanln(&input)
		conn.Write([]byte(input))
		// if s, err := strconv.Atoi(input); err == nil {
		// 	// fmt.Printf("%T, %v", s, s-1)
		// 	// g = g.playCard(true, s-1)

		// }
	}
	conn.Close()
}
