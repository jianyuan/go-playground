package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func signals() <-chan struct{} {
	quit := make(chan struct{})

	go func() {
		signals := make(chan os.Signal)
		defer close(signals)

		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
		defer signal.Stop(signals)

		<-signals
		quit <- struct{}{}
	}()

	return quit
}

func main() {
	quit := signals()

	done := make(chan struct{}, 1)

	go func() {
		log.Println("waiting")
		<-quit
		log.Println("received quit signal")
		done <- struct{}{}
	}()

	<-done
}
