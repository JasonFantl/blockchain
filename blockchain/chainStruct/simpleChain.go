package chainStruct

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type SimpleChain struct {
	Nodes []Node
}

func (sc SimpleChain) Verify() bool {
	for i := 1; i < len(sc.Nodes); i++ {
		if bytes.Compare(sc.Nodes[i].PrevHash, hashNode(sc.Nodes[i-1])) != 0 {
			return false
		}
		if !verifyNode(sc.Nodes[i]) {
			return false
		}
	}
	return true
}

func (sc *SimpleChain) Append(data []byte) {
	if sc.Nodes == nil {
		sc.Nodes = make([]Node, 0)
	}

	var prevHash []byte
	if len(sc.Nodes) > 0 {
		prevHash = hashNode(sc.Nodes[len(sc.Nodes)-1])
	}

	node := Node{
		Data:     data,
		PrevHash: prevHash,
		Nonce:    0,
	}

	// mine the node
	for !verifyNode(node) {
		node.Nonce++
	}

	sc.Nodes = append(sc.Nodes, node)
}

func (sc SimpleChain) ToBytes() ([]byte, error) {
	return json.Marshal(sc)
}

func (sc *SimpleChain) Update(data []byte) error {
	newSC := SimpleChain{}
	err := json.Unmarshal(data, &newSC)
	if err != nil {
		return err
	}

	if len(newSC.Nodes) > len(sc.Nodes) {
		if newSC.Verify() {
			sc.Nodes = newSC.Nodes
			fmt.Println("got updated chain")
		}
	}

	return nil
}

func (sc SimpleChain) Array() [][]byte {

	data := make([][]byte, 0)
	for _, node := range sc.Nodes {
		data = append(data, node.Data)
	}

	return data
}

func (sc SimpleChain) String() string {
	s := fmt.Sprintf("Simple Chain:\n")
	s += fmt.Sprintf("  Nodes:\n")
	for _, node := range sc.Nodes {
		s += node.Stringify(4) + "\n"
	}
	return s
}
