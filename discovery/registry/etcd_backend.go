package registry

import (
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
)

// Warning: ttl and heartbeatPeriod are not thread safe
type EtcdBackend struct {
	client         *etcd.Client
	wg             sync.WaitGroup
	ttl            uint64
	ttlPeriod      time.Duration
	deregisterChan chan struct{}
}

// Interface check trick
var _ Backend = (*EtcdBackend)(nil)

func NewEtcdBackend(machines []string, ttl uint64) Backend {
	return &EtcdBackend{
		client:         etcd.NewClient(machines),
		ttl:            ttl,
		ttlPeriod:      time.Duration(ttl*9/10) * time.Second,
		deregisterChan: make(chan struct{}, 1),
	}
}

func (b *EtcdBackend) Register(n *Node) error {
	key, value := keyFromNode(n), ""

	if err := b.register(key, value); err != nil {
		return err
	}

	b.wg.Add(1)
	go b.registerKeepAlive(key, value)

	return nil
}

func (b *EtcdBackend) Deregister() <-chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		b.wg.Wait()
		close(done)
	}()
	close(b.deregisterChan)
	return done
}

func (b *EtcdBackend) registerKeepAlive(key, value string) {
	defer b.wg.Done()

	logWithNode := log.WithFields(log.Fields{"key": key, "value": value})

	ticker := time.NewTicker(b.ttlPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logWithNode.Debug("Heartbeat")
			b.register(key, value)
		case <-b.deregisterChan:
			logWithNode.Debug("Stopping heartbeat")
			b.deregister(key)
			return
		}
	}
}

func (b *EtcdBackend) register(key, value string) error {
	log.WithFields(log.Fields{"key": key, "value": value}).Debug("Set")
	_, err := b.client.Set(key, value, b.ttl)
	return err
}

func (b *EtcdBackend) deregister(key string) error {
	log.WithField("key", key).Debug("Delete")
	_, err := b.client.Delete(key, false)
	return err
}

func keyFromNode(n *Node) string {
	return fmt.Sprintf("%s/%d/%d/%d/%s", n.Name, n.Version.Major, n.Version.Minor, n.Version.Patch, n.Id)
}
