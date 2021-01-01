package main

import (
	"github.com/Allenxuxu/stark/util/log"
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetPrefix("[test lib]")

	log.Info("hello")
	log.Debug("world")
}
