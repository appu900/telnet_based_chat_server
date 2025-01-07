package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients     = make(map[*Client]bool)
	clientsLock sync.Mutex
)

func main() {
	listner, err := net.Listen("tcp", ":2020")
	if err != nil {
		log.Fatalf("failed to start chat server: %v", err)
	}
	defer listner.Close()

	fmt.Println("Chat server started on :2020")

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	client := &Client{
		conn: conn,
	}
	client.conn.Write([]byte("Welcome to the chat server ! Please enter your name\n"))
	client.conn.Write([]byte("here are the people avalavil in the chat section\n"))
	for client := range clients {
		client.conn.Write([]byte(client.name + "\n"))
	}
	name, err := bufio.NewReader(client.conn).ReadString('\n')
	if err != nil {
		log.Printf("failed to read name: %v", err)
		return
	}
	client.name = strings.TrimSpace(name)
	client.conn.Write([]byte(fmt.Sprintf("Welcome %s !\n", client.name)))

	clientsLock.Lock()
	clients[client] = true
	clientsLock.Unlock()

	brodCast(fmt.Sprintf("%s joined the chat", client.name), nil)
	client.conn.Write([]byte(fmt.Sprintf("Hello,%s! You are now connected to the chat server\n", client.name)))

	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "exit" {
			break
		}
		brodCast(fmt.Sprintf("%s: %s", client.name, message), client)
	}

	// handle client disconnection from server
	clientsLock.Lock()
	delete(clients, client)
	clientsLock.Unlock()
	brodCast(fmt.Sprintf("%s left the chat", client.name), nil)
	log.Printf("client %s disconnected", client.name)

}

func brodCast(message string, sender *Client) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	for client := range clients {
		if client != sender {
			client.conn.Write([]byte(message + "\n"))
		}
	}
}
