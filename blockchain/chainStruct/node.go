package chainStruct

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Node struct {
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func verifyNode(node Node) bool {
	return verifyHash(hashNode(node))
}

func verifyHash(hash []byte) bool {
	difficulty := 1

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

func hashNode(node Node) []byte {
	bytes, err := getBytes(node)
	if err != nil {
		return nil
	}

	return hash(bytes)
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func hash(bytes []byte) []byte {
	hashed := sha256.Sum256(bytes)
	return hashed[:]
}

func (node Node) Stringify(i int) string {
	s := fmt.Sprintf("%*sNode: %x\n", i, "", hashNode(node))
	s += fmt.Sprintf("%*sPrev: %x\n", i+2, "", node.PrevHash)
	s += fmt.Sprintf("%*sData: %s\n", i+2, "", node.Data)

	return s
}
