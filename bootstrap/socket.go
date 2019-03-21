package bootstrap

import (
  "github.com/graarh/golang-socketio"
  "github.com/graarh/golang-socketio/transport"
  "log"
)

const (
  T_INFO 			  = "T_INFO"
  T_WARNING 		= "T_WARNING"
	T_ERROR 		  = "T_ERROR"
	SOCKET_SERVER = "wss://ods-visualbox.cs.uit.no:1337/socket.io/?EIO=3&transport=websocket"
)

type Message struct {
  Type    int    `json:"type"`
  Channel string `json:"channel"`
  Text    string `json:"text"`
}

func onMessage(args Message) {
  switch Message.type {
  case "TICK":
    tick()
  case "TERMINATE":
    i := int32(m.i)

    if i < 0:
      // terminate
  default:
    log.Printf("Unknown socket message type: %v", Message.type)
  }
}

func initSocket() {
	// Setup socket
  c, err := gosocketio.Dial(SOCKET_SERVER, transport.GetDefaultWebsocketTransport())
  if err != nil {
    log.Fatal(err)
  }

  // Send join on socket connection
  err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
    // Join and send init
  })
  if err != nil {
    log.Fatal(err)
  }

  // Abort if desconnected
  err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
    log.Fatal("Disconnected")
  })
  if err != nil {
    log.Fatal(err)
  }

  // Dispatch message handler on socket message
  err = c.On("message", func(h *gosocketio.Channel, args Message) {
    go onMessage(args)
  })
  if err != nil {
    log.Fatal(err)
  }
}
