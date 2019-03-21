package bootstrap

import (
  "log"
  "sync"
)

var (
  wg := &sync.WaitGroup{}
)

func main() {
  // Init
  wg.Add(1)
  go drain()
  initSocket()

  // Standby
  wg.Wait()
}
