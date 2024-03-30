package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"
)

// Constants defining server configurations
const (
	Port        = "6969"
	SafeMode    = true
	MessageRate = 1.0
	BanLimit    = 10 * 60.0
	StrikeLimit = 10
)

// MessageType enumerates the types of messages
type MessageType int

// Enumerated values for MessageType
const (
	ClientConnected MessageType = iota + 1
	ClientDisconnected
	NewMessage
)

// Message represents a communication message
type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

// Client represents a connected client
type Client struct {
	Conn        net.Conn
	LastMessage time.Time
	StrikeCount int
}

// sensitive replaces sensitive information based on safe mode setting
func sensitive(message string) string {
	if SafeMode {
		return "[REDACTED]"
	}
	return message
}

// server listens for messages and handles client interactions
func server(messages chan Message) {
	clients := make(map[string]*Client)
	bannedClients := make(map[string]time.Time)

	for msg := range messages {
		switch msg.Type {
		case ClientConnected:
			handleClientConnected(msg, clients, bannedClients)
		case ClientDisconnected:
			handleClientDisconnected(msg, clients)
		case NewMessage:
			handleNewMessage(msg, clients, bannedClients)
		}
	}
}

// handleClientConnected handles a new client connection
func handleClientConnected(msg Message, clients map[string]*Client, bannedClients map[string]time.Time) {
	addr := msg.Conn.RemoteAddr().(*net.TCPAddr)
	bannedAt, banned := bannedClients[addr.IP.String()]
	now := time.Now()

	if banned && now.Sub(bannedAt).Seconds() < BanLimit {
		msg.Conn.Write([]byte(fmt.Sprintf("You are banned MF: %f secs left\n", BanLimit-now.Sub(bannedAt).Seconds())))
		msg.Conn.Close()
		return
	}

	log.Printf("Client %s connected", sensitive(addr.String()))
	clients[addr.String()] = &Client{
		Conn:        msg.Conn,
		LastMessage: now,
	}
}

// handleClientDisconnected handles a client disconnection
func handleClientDisconnected(msg Message, clients map[string]*Client) {
	addr := msg.Conn.RemoteAddr().(*net.TCPAddr)
	log.Printf("Client %s disconnected", sensitive(addr.String()))
	delete(clients, addr.String())
}

// handleNewMessage handles a new message received from a client
func handleNewMessage(msg Message, clients map[string]*Client, bannedClients map[string]time.Time) {
	authorAddr := msg.Conn.RemoteAddr().(*net.TCPAddr)
	author := clients[authorAddr.String()]
	now := time.Now()

	if author == nil {
		msg.Conn.Close()
		return
	}

	if !isMessageValid(msg.Text) || now.Sub(author.LastMessage).Seconds() < MessageRate {
		author.StrikeCount++
		if author.StrikeCount >= StrikeLimit {
			bannedClients[authorAddr.IP.String()] = now
			author.Conn.Write([]byte("You are banned MF\n"))
			author.Conn.Close()
		}
	} else {
		author.LastMessage = now
		author.StrikeCount = 0
		log.Printf("Client %s sent message %s", sensitive(authorAddr.String()), msg.Text)
		broadcastMessage(msg, clients, authorAddr)
	}
}

// isMessageValid checks if a message contains valid UTF-8 text
func isMessageValid(text string) bool {
	return utf8.ValidString(text)
}

// broadcastMessage broadcasts a message to all clients except the sender
func broadcastMessage(msg Message, clients map[string]*Client, authorAddr *net.TCPAddr) {
	for _, client := range clients {
		if client.Conn.RemoteAddr().String() != authorAddr.String() {
			client.Conn.Write([]byte(msg.Text))
		}
	}
}

// client manages communication with a connected client
func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 64)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			messages <- Message{
				Type: ClientDisconnected,
				Conn: conn,
			}
			return
		}
		text := string(buffer[:n])
		messages <- Message{
			Type: NewMessage,
			Text: text,
			Conn: conn,
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s: %s\n", Port, sensitive(err.Error()))
	}
	defer ln.Close()

	log.Printf("Listening to TCP connections on port %s ...\n", Port)

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept a connection: %s\n", sensitive(err.Error()))
			continue
		}
		messages <- Message{
			Type: ClientConnected,
			Conn: conn,
		}
		go client(conn, messages)
	}
}
