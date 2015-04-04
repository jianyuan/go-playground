package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"
)

func main() {
	b := backoff.NewExponentialBackOff()
	ticker := backoff.NewTicker(b)

	operation := func() error {
		return fmt.Errorf("Errored")
	}

	for t := range ticker.C {
		log.WithFields(log.Fields{"ticker": t}).Info()

		if err := operation(); err != nil {
			log.WithFields(log.Fields{"error": err}).Info("Will retry...")
			continue
		}

		ticker.Stop()
		break
	}
}
