package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

func reverseCards(input []Card) []Card {
	if len(input) == 0 {
		return input
	}
	return append(reverseCards(input[1:]), input[0])
}

func playCardAutomaticaly(playerInputChannel chan int) {
	time.Sleep(1 * time.Second)
	playerInputChannel <- 1
}

func updateGui(gameStateChannel chan NetworkMessage, playerInputChannel chan int, opponentButtons [4]*widget.Button, buttons [4]*widget.Button, boardButtons [4]*widget.Button, w fyne.Window) {
	for {
		nm := <-gameStateChannel

		if nm.GameOverState.IsGameOver {
			sPlayerCards := []string{"Player Won Cards (" + fmt.Sprint(len(nm.GameOverState.PlayerWonCards)) + "): "}
			sOpponentCards := []string{"Opponent Won Cards (" + fmt.Sprint(len(nm.GameOverState.OpponentWonCards)) + "): "}
			for i := 0; i < len(nm.GameOverState.PlayerWonCards); i++ {
				sPlayerCards = append(sPlayerCards, nm.GameOverState.PlayerWonCards[i].toString())
				if i%10 == 9 {
					sPlayerCards = append(sPlayerCards, "\n")
				}
			}
			for i := 0; i < len(nm.GameOverState.OpponentWonCards); i++ {
				sOpponentCards = append(sOpponentCards, nm.GameOverState.OpponentWonCards[i].toString())
				if i%10 == 9 {
					sOpponentCards = append(sOpponentCards, "\n")
				}
			}
			// tPlayerCards := container.NewHBox(widget.NewList(
			// 	func() int {
			// 		return len(nm.GameOverState.PlayerWonCards)
			// 	},
			// 	func() fyne.CanvasObject {
			// 		return widget.NewLabel("template")
			// 	},
			// 	func(i widget.ListItemID, o fyne.CanvasObject) {
			// 		o.(*widget.Label).SetText(nm.GameOverState.PlayerWonCards[i].toString())
			// 	}))
			// tPlayerCards.Resize(fyne.NewSize(500, 60))

			// tOpponentCards := container.NewHBox(widget.NewList(
			// 	func() int {
			// 		return len(nm.GameOverState.OpponentWonCards)
			// 	},
			// 	func() fyne.CanvasObject {
			// 		return widget.NewLabel("template")
			// 	},
			// 	func(i widget.ListItemID, o fyne.CanvasObject) {
			// 		o.(*widget.Label).SetText(nm.GameOverState.OpponentWonCards[i].toString())
			// 	}))
			// tOpponentCards.Resize(fyne.NewSize(500, 60))
			playerLabel := widget.NewLabel(strings.Join(sPlayerCards, "-"))
			opponentLabel := widget.NewLabel(strings.Join(sOpponentCards, "-"))
			playerLabel.Resize(fyne.NewSize(400, 50))
			opponentLabel.Resize(fyne.NewSize(400, 50))

			tVbox := container.NewVBox(playerLabel, opponentLabel)
			tVbox.Resize(fyne.NewSize(500, 150))

			d := dialog.NewCustom("Game Over", "Finish", tVbox, w)
			d.Resize(fyne.NewSize(700, 200))
			d.Show()
		} else {
			// fmt.Printf("New Update State: %s\n", fmt.Sprint(nm))
			for i := 0; i < 4; i++ {
				if i < len(nm.PlayerHand) {
					buttons[i].Show()
					buttons[i].SetText(nm.PlayerHand[i].toString())
					if nm.IsPlayerTurn {
						buttons[i].Enable()
						go playCardAutomaticaly(playerInputChannel)
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

			reversedOpenCards := reverseCards(nm.BoardOpenCards)
			fmt.Println(nm.BoardOpenCards)
			fmt.Println(reversedOpenCards)
			for i := 0; i < 4; i++ {
				if i < len(reversedOpenCards) {
					boardButtons[i].SetText(reversedOpenCards[i].toString())
				} else {
					boardButtons[i].SetText(UNKNOWN)
				}
				boardButtons[i].Refresh()
			}
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
	w.Resize(fyne.Size{Width: 950, Height: 450})
	w.SetContent(parentContainer)

	go updateGui(gameStateChannel, playerInputChannel, opponentButtons, buttons, boardButtons, w)
	w.ShowAndRun()
}

func runClientGui() {
	gameStateChannel := make(chan NetworkMessage)
	playerInputChannel := make(chan int)
	go runClient(gameStateChannel, playerInputChannel)
	showGui(gameStateChannel, playerInputChannel)
}
