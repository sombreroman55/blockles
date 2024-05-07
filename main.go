package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/sombreroman55/blockles/site"
)

func main() {
	log.Printf("Hello World!")
	site.ServeBlockles()
}
