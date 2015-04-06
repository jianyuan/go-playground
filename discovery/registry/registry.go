package registry

type Registry struct {
	backend Backend
}

func NewRegistry(backend Backend) *Registry {
	return &Registry{
		backend: backend,
	}
}

func (reg *Registry) Register(n *Node) {
	reg.backend.Register(n)
}

func (reg *Registry) Deregister() {
	<-reg.backend.Deregister()
}
