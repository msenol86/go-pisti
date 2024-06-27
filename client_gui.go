package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func handleServerConnection(c net.Conn, gameStateChannel chan NetworkMessage) {
	for {
		d := json.NewDecoder(c)

		for {
			var msg NetworkMessage
			if err := d.Decode(&msg); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			gameStateChannel <- msg

		}
	}
	c.Close()
}

func handlePlayerInput(c net.Conn, playerInputChannel chan int) {
	for {
		input := <-playerInputChannel
		c.Write([]byte(PLAY + " " + strconv.Itoa(input) + "\n"))
	}
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
	go handleServerConnection(conn, gameStateChannel)
	go handlePlayerInput(conn, playerInputChannel)
}

func updateGui(gameStateChannel chan NetworkMessage, opponentButtons [4]*widget.Button, buttons [4]*widget.Button, boardButtons [4]*widget.Button) {
	for {
		nm := <-gameStateChannel
		// fmt.Printf("New Update State: %s\n", fmt.Sprint(nm))
		for i := 0; i < 4; i++ {
			if i < len(nm.PlayerHand) {
				buttons[i].Show()
				buttons[i].SetText(fmt.Sprint(nm.PlayerHand[i]))
				if nm.IsPlayerTurn {
					buttons[i].Enable()
				} else {
					buttons[i].Disable()
				}

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

func runClientGui() {
	gameStateChannel := make(chan NetworkMessage)
	playerInputChannel := make(chan int)
	go runClient(gameStateChannel, playerInputChannel)
	showGui(gameStateChannel, playerInputChannel)
}