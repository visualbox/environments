package main

import (
	"log"

	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/mtfelian/golang-socketio/transport"
)

const (
	// StatusTypeInfo ...
	StatusTypeInfo = "T_INFO"
	// StatusTypeWarning ...
	StatusTypeWarning = "T_WARNING"
	// StatusTypeError ...
	StatusTypeError = "T_ERROR"

	messageTypeInit   = "INIT"
	messageTypeStatus = "STATUS"
	socketServer      = "localhost" // ods-visualbox.cs.uit.no
	sockerPort        = 1337
)

type message struct {
	Type       string `json:"type"`
	I          string `json:"i,omitempty"`
	StatusType string `json:"statusType,omitempty"`
	Data       string `json:"data,omitempty"`
}

var (
	socketChannel *gosocketio.Channel
)

// Status ...
func Status(statusType string, data string) {
	if socketChannel == nil {
		return
	}

	message := message{
		Type:       messageTypeStatus,
		StatusType: statusType,
		Data:       data,
	}
	log.Printf("Sending status: %v", message)
	err := socketChannel.Emit("message", message)
	if err != nil {
		log.Println(err)
	}
}

func onConnectionHandler(c *gosocketio.Channel) {
	// Join
	err := c.Emit("join", EnvToken)
	if err != nil {
		log.Fatal(err)
	}
	// Send INIT
	err = c.Emit("message", message{
		Type: messageTypeInit,
		I:    EnvI,
	})
	if err != nil {
		log.Fatal(err)
	}

	socketChannel = c
	wg.Done()
}

func onDisconnectionHandler(c *gosocketio.Channel) {
	log.Fatal("Disconnected")
}

func onMessageHandler(c *gosocketio.Channel, data interface{}) {
	log.Printf("--- Client channel %s received someEvent with data: %v\n", c.Id(), data)
	/*j, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}*/

	/*
			// OnMessage ...
		func OnMessage(args Message) {
			switch args.Type {
			case "TICK":
				Tick()
			default:
				log.Printf("Unknown socket message type: %v\n", args.Type)
			}
		}
	*/
}

// InitSocket ...
func InitSocket() {
	// Setup socket
	client, err := gosocketio.Dial(
		gosocketio.AddrWebsocket(socketServer, sockerPort, false),
		transport.DefaultWebsocketTransport(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.On(gosocketio.OnConnection, onConnectionHandler); err != nil {
		log.Fatal(err)
	}
	if err := client.On(gosocketio.OnDisconnection, onDisconnectionHandler); err != nil {
		log.Fatal(err)
	}

	if err := client.On("message", onMessageHandler); err != nil {
		log.Fatal(err)
	}
}
