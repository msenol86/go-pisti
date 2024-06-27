package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

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
const UNKNOWN = "� �"

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

func updateGui(gameStateChannel chan NetworkMessage, opponentButtons [4]*widget.Button, buttons [4]*widget.Button, boardButtons [4]*widget.Button) {
	for {
		nm := <-gameStateChannel
		for i := 0; i < 4; i++ {
			if i < len(nm.PlayerHand) {
				buttons[i].Show()
				buttons[i].SetText(fmt.Sprint(nm.PlayerHand[i]))
				buttons[i].Enable()
			} else {
				buttons[i].Hide()
				buttons[i].Disable()
			}
			buttons[i].Refresh()
		}
		for i := 0; i < 4; i++ {
			if i < int(nm.OpponentHandCount) {
				opponentButtons[i].Show()
			} else {
				opponentButtons[i].Hide()
			}
			opponentButtons[i].Refresh()
		}

		for i := 0; i < min(4, int(nm.BoardCount)); i++ {
			if i < len(nm.BoardOpenCards) {
				boardButtons[i].Show()
				boardButtons[i].SetText(fmt.Sprint(nm.BoardOpenCards[i]))
			} else {
				boardButtons[i].SetText(UNKNOWN)
			}
			boardButtons[i].Refresh()
		}
	}
}

func showGui(gameStateChannel chan NetworkMessage, playerInputChannel chan int) {
	a := app.New()
	w := a.NewWindow("Pisti")

	opponentButtons := [4]*widget.Button{
		widget.NewButton(UNKNOWN, func() {

		}),
		widget.NewButton(UNKNOWN, func() {

		}),
		widget.NewButton(UNKNOWN, func() {

		}),
		widget.NewButton(UNKNOWN, func() {

		})}
	buttons := [4]*widget.Button{
		widget.NewButton(UNKNOWN, func() {
			playerInputChannel <- 1
		}),
		widget.NewButton(UNKNOWN, func() {
			playerInputChannel <- 2
		}),
		widget.NewButton(UNKNOWN, func() {
			playerInputChannel <- 3
		}),
		widget.NewButton(UNKNOWN, func() {
			playerInputChannel <- 4
		})}

	boardButtons := [4]*widget.Button{
		widget.NewButton(UNKNOWN, func() {

		}),
		widget.NewButton(UNKNOWN, func() {

		}),
		widget.NewButton(UNKNOWN, func() {

		}),
		widget.NewButton(UNKNOWN, func() {

		})}

	for i := 0; i < len(buttons); i++ {
		buttons[i].Disable()
		boardButtons[i].Disable()
		opponentButtons[i].Disable()
	}
	topContainer := container.NewHBox(
		opponentButtons[0],
		opponentButtons[1],
		opponentButtons[2],
		opponentButtons[3])
	boardContainer := container.NewHBox(
		boardButtons[0],
		boardButtons[1],
		boardButtons[2],
		boardButtons[3])
	bottomContainer := container.NewHBox(
		buttons[0],
		buttons[1],
		buttons[2],
		buttons[3],
	)
	parentContainer := container.NewVBox(
		topContainer,
		boardContainer,
		bottomContainer,
	)
	// w.Resize(fyne.Size{Width: 800, Height: 800})
	w.SetContent(parentContainer)

	go updateGui(gameStateChannel, opponentButtons, buttons, boardButtons)
	w.ShowAndRun()
}

func main() {
	gameStateChannel := make(chan NetworkMessage)
	playerInputChannel := make(chan int)
	go runClient(gameStateChannel, playerInputChannel)
	showGui(gameStateChannel, playerInputChannel)
}
