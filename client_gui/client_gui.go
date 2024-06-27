package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	HOST = "localhost"
	PORT = "6666"
	TYPE = "tcp"
	JOIN = "JOIN"
	PLAY = "PLAY"
)

type SuitType string
type RankType int

type Card struct {
	Suit SuitType
	Rank RankType
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
}

func runClient(gameStateChannel chan NetworkMessage, playerInputChannel chan int) {
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
	_, err = conn.Write([]byte(JOIN + "\n"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}
	for {
		// buffer to get data
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		var nm NetworkMessage
		d := json.NewDecoder(conn)
		if err := d.Decode(&nm); err != nil {
			fmt.Println(err)
		}
		// boardArray := []string{}
		// for i := 0; i < int(nm.BoardCount); i++ {
		// 	if i >= int(nm.BoardCount)-len(nm.BoardOpenCards) {
		// 		theCard := nm.BoardOpenCards[i-(int(nm.BoardCount)-len(nm.BoardOpenCards))]
		// 		boardArray = append(boardArray, fmt.Sprint(theCard))
		// 	} else {
		// 		boardArray = append(boardArray, "{* *}")
		// 	}
		// }
		gameStateChannel <- nm
		// fmt.Print("Won Cards Count: ", nm.PlayerWonCardsCount)
		// fmt.Print(" - Points: ", nm.PlayerPoints)
		// fmt.Println(" - Pist Count: ", nm.PlayerPistiCounts)
		// fmt.Println("Board: ", boardArray)
		// fmt.Println("Hand", nm.PlayerHand)

		// fmt.Printf("Type the number of card you want to play [%d-%d]: ", 1, len(nm.PlayerHand))
		// time.Sleep(2 * time.Second)
		// input := "1"
		// var input string
		// fmt.Scanln(&input)
		input := <-playerInputChannel
		conn.Write([]byte(PLAY + " " + strconv.Itoa(input) + "\n"))
	}
	conn.Close()
}

func showGui(gameStateChannel chan NetworkMessage, playerInputChannel chan int) {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	buttonTexts := [4]string{"* *", "* *", "* *", "* *"}
	buttons := [4]*widget.Button{
		widget.NewButton(buttonTexts[0], func() {
			playerInputChannel <- 1
		}),
		widget.NewButton(buttonTexts[1], func() {
			playerInputChannel <- 2
		}),
		widget.NewButton(buttonTexts[2], func() {
			playerInputChannel <- 3
		}),
		widget.NewButton(buttonTexts[3], func() {
			playerInputChannel <- 4
		})}
	w.SetContent(container.NewHBox(
		hello,
		buttons[0],
		buttons[1],
		buttons[2],
		buttons[3],
	))
	// str := binding.NewString()
	// str.Set("Initial value")

	// text := widget.NewLabelWithData(str)
	// w.SetContent(text)

	counter := 0
	go func() {
		for {
			time.Sleep(time.Second * 2)
			// str.Set("A new string " + fmt.Sprint(counter))
			// buttonTexts[0] = "TEST"
			buttons[counter].SetText("Test")
			buttons[counter].Refresh()
			counter += 1
			if counter > 4 {
				break
			}
		}
	}()

	w.ShowAndRun()
}

func main() {
	gameStateChannel := make(chan NetworkMessage)
	playerInputChannel := make(chan int)
	// go runClient(gameStateChannel, playerInputChannel)
	showGui(gameStateChannel, playerInputChannel)
}
