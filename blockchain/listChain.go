package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ListChain struct {
	Nodes map[string][]Block
}

func (lc ListChain) Verify() bool {
	for _, txs := range lc.Nodes {
		for i := 1; i < len(txs); i++ {
			if bytes.Compare(txs[i].Signed.PrevHash, hash(txs[i-1])) != 0 {
				return false
			}
			if !verifyBlock(txs[i]) {
				return false
			}
		}
	}
	return true
}

func (lc *ListChain) Append(tx Transaction, wallet Wallet) error {

	if lc.Nodes == nil {
		lc.Nodes = make(map[string][]Block)
	}

	txs, ok := lc.Nodes[string(tx.From)]

	if !ok {
		lc.Nodes[string(tx.From)] = make([]Block, 0)
		txs = lc.Nodes[string(tx.From)]
	}

	var prevHash []byte
	if len(txs) > 0 {
		prevHash = hash(txs[len(txs)-1])
	}

	node := Node{
		Transaction: tx,
		PrevHash:    prevHash,
	}

	sig, err := wallet.signMessage(node)
	if err != nil {
		return err
	}

	block := Block{
		Signature: sig,
		Nonce:     0,
		Signed:    node,
	}

	// mine the node
	for !verifyHash(hash(block)) {
		block.Nonce++
	}

	lc.Nodes[string(tx.From)] = append(txs, block)
	return nil
}

func (lc ListChain) ToBytes() ([]byte, error) {
	jsoned, err := json.Marshal(lc)
	return jsoned, err
}

func (lc *ListChain) Update(data []byte) error {

	newSC := ListChain{}
	err := json.Unmarshal(data, &newSC)
	if err != nil {
		return err
	}

	if newSC.Verify() {
		lc.Nodes = newSC.Nodes
		fmt.Println("got updated chain")
	}

	return nil
}

func (lc ListChain) String() string {
	s := fmt.Sprintf("List Chain:\n")

	for _, nodes := range lc.Nodes {
		for _, node := range nodes {
			s += node.Stringify(4) + "\n"
		}
		s += fmt.Sprintln()
	}
	return s
}
