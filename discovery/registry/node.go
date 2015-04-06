package registry

type Node struct {
	Name    string
	Version Version
	Id      string
}

func NewNode(name string, version Version, id string) *Node {
	return &Node{name, version, id}
}
