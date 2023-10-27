package model

type Node struct {
	id        string
	memory    float64
	coreCount float64
}

func NewNode(id string, memory float64, coreCount float64) *Node {
	return &Node{
		id:        id,
		memory:    memory,
		coreCount: coreCount,
	}
}

func (n *Node) Id() string {
	return n.id
}

func (n *Node) Memory() float64 {
	return n.memory
}

func (n *Node) Cores() float64 {
	return n.coreCount
}
