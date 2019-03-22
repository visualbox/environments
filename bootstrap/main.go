package main

import (
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
	// Init
	wg.Add(1)
	go drain()
	InitSocket()
	StartIntegration()

	// Standby
	wg.Wait()
}
