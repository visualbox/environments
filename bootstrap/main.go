package main

import (
	"log"
	"os"
	"sync"
)

var (
	// EnvI - Integration dashboard index.
	EnvI = os.Getenv("I")

	// EnvToken - Instance token.
	EnvToken = os.Getenv("TOKEN")

	// EnvRestAPIID - AWS API GW ID
	EnvRestAPIID = os.Getenv("REST_API_ID")

	// EnvID - VisualBox integration ID
	EnvID = os.Getenv("ID")

	// EnvVersion - Integration version.
	EnvVersion = os.Getenv("VERSION")

	// EnvModel - Initial integration configuration model.
	EnvModel = os.Getenv("MODEL")

	wg = &sync.WaitGroup{}
)

func main() {
	log.Println("Init Unix socket server")
	go InitUnixSocket()

	log.Println("Init Socket.io client")
	wg.Add(1)
	InitSocket()
	wg.Wait()

	log.Println("Init project files and start drainage")
	wg.Add(1)
	go StartIntegration()
	go Drain()
	wg.Wait()

	Terminate(true)
}
