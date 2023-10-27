package model

type Node struct {
	id        int
	memory    float64
	coreCount float64
}

func NewNode(id int, memory float64, coreCount float64) *Node {
	return &Node{
		id:        id,
		memory:    memory,
		coreCount: coreCount,
	}
}

func (n *Node) Id() int {
	return n.id
}

func (n *Node) Memory() float64 {
	return n.memory
}

func (n *Node) Cores() float64 {
	return n.coreCount
}
