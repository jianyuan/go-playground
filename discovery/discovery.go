package main

import (
	"code.google.com/p/go-uuid/uuid"
	log "github.com/Sirupsen/logrus"

	"github.com/jianyuan/go-playground/discovery/registry"
	"github.com/jianyuan/go-playground/signals"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	quit := signals.TerminateSignal()

	machines := []string{"http://127.0.0.1:2379", "http://127.0.0.1:4001"}
	backend := registry.NewEtcdBackend(machines, 10)
	reg := registry.NewRegistry(backend)

	instanceId := uuid.NewUUID().String()

	for _, node := range []*registry.Node{
		registry.NewNode("healthcheck", registry.NewVersion(1, 0, 0), instanceId),
		registry.NewNode("signup", registry.NewVersion(1, 0, 0), instanceId),
		registry.NewNode("signin", registry.NewVersion(1, 0, 0), instanceId),
	} {
		reg.Register(node)
	}
	defer func() {
		reg.Deregister()
		log.Info("Deregistered gracefully")
	}()

	<-quit
}
