package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn net.Conn
}

func (client *Client) handleRequest() {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		fmt.Printf("Message incoming: %s", string(message))
		client.conn.Write([]byte("Message received.\n"))
	}
}

func startServer(data chan int) {
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
		go client.handleRequest()
	}
}
