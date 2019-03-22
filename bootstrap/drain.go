package main

import (
	"log"
	"time"
)

const (
	timeout     = 60
	timeoutTick = 5
)

var (
	lastCheck = int32(time.Now().Unix())
)

// Tick -Update last checked timestamp.
func Tick() {
	lastCheck = int32(time.Now().Unix())
}

// Drain - each TIMEOUT_TICK interval
// and check that the diff of lastCheck and
// now isn't over TIMEOUT.
func Drain() {
	for {
		now := int32(time.Now().Unix())
		diff := now - lastCheck

		log.Printf("tick (terminating in %v seconds)\n", timeout-diff)

		if diff >= timeout {
			log.Println("TIMEOUT")
		}

		time.Sleep(timeoutTick * time.Second)
	}
}
