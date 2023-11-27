package core

// Node
// TODO Incomplete implementation.
type Node struct {
	broker Broker
}

type NodeConfig struct {
}

func NewNode(config NodeConfig) *Node {
	n := &Node{}
	return n
}

func (n *Node) SetBroker(b Broker) {
	n.broker = b
}
