package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	conn net.Conn
}

func (client *Client) handleChannelMessages(data chan NetworkMessage) {
	for {
		msg := <-data
		ajson, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("JSON Error", msg)
		}
		fmt.Println(msg)
		client.conn.Write([]byte(ajson))
		// foo_marshalled, err := json.Marshal(Foo{Number: 1, Title: "test"})
		// fmt.Fprint(w, string(foo_marshalled)) // write response to ResponseWriter (w)
	}
}

func (client *Client) handleRequest(joinChannel chan string, inputChannel chan int) {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		msg_str := strings.TrimSpace(string(message))
		fmt.Printf("Message incoming: %s", string(msg_str))
		client.conn.Write([]byte("Message received.\n"))
		if msg_str == "JOIN" {
			// fmt.Printf("\nPlayer joined: %s\n", client.conn.RemoteAddr())
			chan_message := client.conn.RemoteAddr().String()
			joinChannel <- chan_message
		} else if strings.HasPrefix(msg_str, "PLAY") {
			words := strings.Split(msg_str, " ")
			if len(words) > 1 {
				if s, err := strconv.Atoi(words[1]); err == nil {
					inputChannel <- s
				}
			}
		}
	}
}

func startServer(data1 chan NetworkMessage, data2 chan NetworkMessage, joinChannel chan string, input1Channel chan int, input2Channel chan int) {
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	counter := 0
	var clients = [2]*Client{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		if counter < 2 {
			clients[0] = client
			fmt.Println("Connected clients ", counter+1)
			if counter < 1 {
				go client.handleRequest(joinChannel, input1Channel)
				go client.handleChannelMessages(data1)
			} else {
				go client.handleRequest(joinChannel, input2Channel)
				go client.handleChannelMessages(data2)
			}

		}
		counter += 1
	}
}
