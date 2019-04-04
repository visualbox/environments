package visualbox

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"reflect"
)

// Model - Intergation configuration model.
var (
	Model map[string]interface{}
	c     net.Conn
	err   error
)

func init() {
	if err = json.Unmarshal([]byte(os.Getenv("MODEL")), &Model); err != nil {
		panic(err)
	}

	if c, err = net.Dial("unix", "/tmp/out"); err != nil {
		panic(err)
	}
}

// Output - send data back to VisualBox
func Output(message interface{}) error {
	var messageString string

	if reflect.TypeOf(message).String() != "string" {
		jsonMap, err := json.Marshal(message)
		if err != nil {
			return err
		}
		messageString = string(jsonMap)
	} else {
		messageString = message.(string)
	}

	msgBytes := []byte(messageString)
	headerBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(headerBuf, uint32(len(msgBytes)))
	_, err := c.Write(append(headerBuf[:], msgBytes[:]...))
	return err
}
