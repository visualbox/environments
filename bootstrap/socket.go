package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
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
	messageTypeOutput = "OUTPUT"
	socketServer      = "ods-visualbox.cs.uit.no"
	socketPort        = 1337
)

type message struct {
	Type       string `json:"type"`
	I          string `json:"i,omitempty"`
	StatusType string `json:"statusType,omitempty"`
	Data       string `json:"data,omitempty"`
}

type clientMessage struct {
	Type        string `json:"type,omitempty"`
	I           string `json:"i,omitempty"`
	Integration struct {
		I       string `json:"i,omitempty"`
		ID      string `json:"id,omitempty"`
		Version string `json:"version,omitempty"`
		Model   string `json:"model,omitempty"`
	} `json:"integration,omitempty"`
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

// Output ...
func Output(data string) {
	if socketChannel == nil {
		return
	}

	message := message{
		Type: messageTypeOutput,
		I:    EnvI,
		Data: data,
	}
	log.Printf("Sending output: %v", message)
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
		Terminate(true)
	}
	// Send INIT
	err = c.Emit("message", message{
		Type: messageTypeInit,
		I:    EnvI,
	})
	if err != nil {
		log.Fatal(err)
		Terminate(true)
	}

	socketChannel = c
	wg.Done()
}

func onDisconnectionHandler(c *gosocketio.Channel) {
	log.Fatal("Disconnected")
	Terminate(true)
}

func onMessageHandler(c *gosocketio.Channel, data interface{}) {
	log.Printf("--- Client channel %s received someEvent with data: %v\n", c.Id(), data)
	j, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	var result map[string]interface{}
	if err = json.Unmarshal(j, &result); err != nil {
		log.Println(err)
		return
	}

	switch result["type"] {
	case "TICK":
		Tick()
	case "TERMINATE":
		// Kill integration process and container
		// if 'i' is not present or same as EnvI.
		if result["i"] == nil || result["i"] == EnvI {
			Terminate(true)

			// Kill process (but not container) if 'i' is < 0 (or anything else)
		} else {
			Terminate(false)
		}
	case "START":
		integration := result["integration"].(map[string]interface{})

		if integration["i"] != EnvI {
			return
		}

		// Set new ENV based on message.
		// User could have changed version, model etc.
		EnvID = integration["id"].(string)
		EnvVersion = integration["version"].(string)

		j, err := json.Marshal(integration["model"])
		if err != nil {
			EnvModel = "{}"
		}
		EnvModel = string(j)

		go StartIntegration()
	default:
		log.Printf("Unknown socket message type: %v\n", result["type"])
	}
}

// InitSocket ...
func InitSocket() {
	// Setup socket
	url := fmt.Sprintf("wss://%s:%d/socket.io/?EIO=3&transport=websocket", socketServer, socketPort)
	client, err := gosocketio.Dial(
		url,
		transport.NewWebsocketTransport(transport.WebsocketTransportParams{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
		Terminate(true)
	}

	if err := client.On(gosocketio.OnConnection, onConnectionHandler); err != nil {
		log.Fatal(err)
		Terminate(true)
	}
	if err := client.On(gosocketio.OnDisconnection, onDisconnectionHandler); err != nil {
		log.Fatal(err)
		Terminate(true)
	}

	if err := client.On("message", onMessageHandler); err != nil {
		log.Fatal(err)
		Terminate(true)
	}
}
