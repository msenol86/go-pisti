package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	conn net.Conn
}

func (client *Client) handleRequest(data chan string, joinChannel chan string) {
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
			fmt.Printf("\nPlayer joined: %s\n", client.conn.RemoteAddr())
			chan_message := client.conn.RemoteAddr().String()
			joinChannel <- chan_message
		}
	}
}

func startServer(data chan string, joinChannel chan string) {
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest(data, joinChannel)
	}
}
