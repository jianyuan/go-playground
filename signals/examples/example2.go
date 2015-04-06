package main

import (
	"log"

	"github.com/jianyuan/go-playground/signals"
)

func main() {
	quit := make(chan struct{}, 1)
	signals.NotifyTerminate(quit)
	log.Println("waiting")
	<-quit
	log.Println("received quit signal")
}
