package main

import (
	"log"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

// wss://ods-visualbox.cs.uit.no:1337/socket.io/?EIO=3&transport=websocket

const (
	// MessageTypeInfo ...
	MessageTypeInfo = "T_INFO"
	// MessageTypeWarning ...
	MessageTypeWarning = "T_WARNING"
	// MessageTypeError ...
	MessageTypeError = "T_ERROR"
	// SocketServer ...
	SocketServer = "ws://localhost:1337/socket.io/?EIO=3&transport=websocket"
)

// Message ...
type Message struct {
	Type       string `json:"type"`
	I          string `json:"i,omitempty"`
	StatusType string `json:"statusType,omitempty"`
	Data       string `json:"data,omitempty"`
}

var (
	// Client ...
	Client gosocketio.Client
)

// Status ...
func Status(statusType string, data string) {
	message := Message{
		Type:       "STATUS",
		StatusType: statusType,
		Data:       data,
	}
	err := Client.Emit("message", message)
	if err != nil {
		log.Println(err)
	}
}

// OnMessage ...
func OnMessage(args Message) {
	switch args.Type {
	case "TICK":
		tick()
	default:
		log.Printf("Unknown socket message type: %v\n", args.Type)
	}
}

// InitSocket ...
func InitSocket() {
	// Setup socket
	Client, err := gosocketio.Dial(SocketServer, transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	// Send join on socket connection
	err = Client.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		// Join
		err = Client.Emit("join", EnvToken)
		if err != nil {
			log.Fatal(err)
		}

		// Send INIT
		err = Client.Emit("message", Message{
			Type: "INIT",
			I:    EnvI,
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	// Abort if disconnected
	err = Client.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	// Dispatch message handler on socket message
	err = Client.On("message", func(h *gosocketio.Channel, args Message) {
		go OnMessage(args)
	})
	if err != nil {
		log.Println(err)
	}
}
