package main

import (
	"log"
	"sync"
)

var (
	// EnvI - Integration dashboard index.
	EnvI = "_0" // os.Getenv("I")

	// EnvToken - Instance token.
	EnvToken = "340599c44386f8f4c18e8ec4b1576ea8d257ba88a61213fbe66ec1b0ad8b59f6aa08a2595d58ea6fa5d4c2f28649b6eee66405aef41b7dc5e069e9f3ee5caa67" // os.Getenv("TOKEN")

	// EnvRestAPIID - AWS API GW ID
	EnvRestAPIID = "grpactt7f8" // os.Getenv("REST_API_ID")

	// EnvID - VisualBox integration ID
	EnvID = "6ba519b0-490f-11e9-b4e4-cb9c67a0ab3e" // os.Getenv("ID")

	// EnvVersion - Integration version.
	EnvVersion = "^" // os.Getenv("VERSION")

	// EnvModel - Initial integration configuration model.
	EnvModel = "{}" // os.Getenv("MODEL")

	wg = &sync.WaitGroup{}
)

func main() {
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
