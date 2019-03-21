package bootstrap

import (
  "time"
)

const (
  /**
   *  Total amount of seconds until termination.
   */
  TIMEOUT = 60

  /**
   *  Interval in seconds when to check timeout.
   */
  TIMEOUT_TICK = 5
)

var (
  lastCheck := int32(time.Now().Unix())
)

/**
 *  Update last checked timestamp.
 */
func tick() {
  lastCheck = int32(time.Now().Unix())
}

/**
 *  Run drain() each TIMEOUT_TICK interval
 *  and check that the diff of lastCheck and
 *  now isn't over TIMEOUT.
 */
func drain() {
  for {
    now := int32(time.Now().Unix())
    diff := now - lastCheck

    log.Printf("tick (terminating in %v seconds)\n", TIMEOUT - diff)

    if diff >= TIMEOUT {
      log.Println("TIMEOUT")
    }

    time.Sleep(TIMEOUT_TICK * time.Second)
  }
}
