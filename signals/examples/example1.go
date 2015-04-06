package main

import (
	"log"
	"sync"

	"github.com/jianyuan/go-playground/signals"
)

func main() {
	var wg sync.WaitGroup

	quit := signals.TerminateSignal()

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("waiting in goroutine 1")
		<-quit
		log.Println("received quit signal 1")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("waiting in goroutine 2")
		<-quit
		log.Println("received quit signal 2")
	}()

	log.Println("waiting in main")
	wg.Wait()
	log.Println("finished waiting")
}
