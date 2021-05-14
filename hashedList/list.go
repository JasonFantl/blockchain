package hashedList

import "fmt"

type List struct {
	nodes       []Node
	queuedNodes []Node
}

func (list *List) QueueNode(node Node) {
	if list.queuedNodes == nil {
		list.queuedNodes = make([]Node, 0)
	}
	list.queuedNodes = append(list.queuedNodes, node)
}

func (list *List) Mine() {
	for _, queuedNode := range list.queuedNodes {

		// set prev hash
		if len(list.nodes) > 0 {
			queuedNode.PrevHash = hashNode(list.nodes[len(list.nodes)-1])
		} else {
			queuedNode.PrevHash = ""
		}

		// mine the block
		for !verifyNode(queuedNode) {
			queuedNode.Nonce++
		}

		// add it into the chain
		list.addNode(queuedNode)
	}

	// clear queue
	list.queuedNodes = nil
}

func (list *List) addNode(block Node) {
	if list.nodes == nil {
		list.nodes = make([]Node, 0)
	}

	if verifyNode(block) {
		list.nodes = append(list.nodes, block)
	} else {
		fmt.Println("invalid node, will not add to list")
	}
}

func (list List) Verify() bool {
	for i := range list.nodes {
		if i != 0 && list.nodes[i].PrevHash != hashNode(list.nodes[i-1]) {
			return false
		}
		if !verifyNode(list.nodes[i]) {
			return false
		}
	}
	return true
}
