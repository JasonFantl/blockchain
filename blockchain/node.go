package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Block struct {
	Nonce     int
	Signature []byte
	Signed    Node
}

type Node struct {
	Transaction Transaction
	PrevHash    []byte
}

func verifyBlock(block Block) bool {
	return verifyHash(hash(block)) && checkSignature(block.Signed, block.Signature, block.Signed.Transaction.From)
}

func verifyHash(hash []byte) bool {
	difficulty := 2

	if len(hash) < difficulty {
		return false
	}

	for i := 0; i < difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

func hash(key interface{}) []byte {
	bytes, err := json.Marshal(key)
	if err != nil {
		return nil
	}
	hashed := sha256.Sum256(bytes)

	return hashed[:]
}

func (block Block) Stringify(i int) string {
	s := fmt.Sprintf("%*sBlock: %x\n", i, "", hash(block))
	s += fmt.Sprintf("%*sTransaction: %s\n", i+2, "", block.Signed.Transaction.String())
	s += fmt.Sprintf("%*sSignature: %x\n", i+2, "", Shorten(block.Signature))

	s += fmt.Sprintf("%*sPrev: %x\n", i+2, "", block.Signed.PrevHash)
	s += fmt.Sprintf("%*sNonce: %d\n", i+2, "", block.Nonce)

	return s
}
