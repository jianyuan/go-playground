package registry

type Backend interface {
	Register(*Node) error
	Deregister() <-chan struct{}
}
