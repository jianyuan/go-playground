package signals

import (
	"os"
	"os/signal"
	"syscall"
)

// TerminateSignal returns a channel that closes when SIGQUIT, SIGTERM or SIGINT signal is received
func TerminateSignal() <-chan struct{} {
	quit := make(chan struct{})
	NotifyTerminate(quit)
	return quit
}

// NotifyTerminate closes the input quit channel when SIGQUIT, SIGTERM or SIGINT signal is received
func NotifyTerminate(quit chan<- struct{}) {
	go func() {
		signals := make(chan os.Signal, 1)
		defer close(signals)

		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
		defer signal.Stop(signals)

		<-signals
		close(quit)
	}()
}
