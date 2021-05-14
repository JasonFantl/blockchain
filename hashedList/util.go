package hashedList

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func hash(bytes []byte) string {
	hashed := sha256.Sum256(bytes)
	return string(hashed[:])
}

func SaveToJSON(list List, path string) error {
	file, err := json.MarshalIndent(list, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, file, 0644)
}

func ImportFromJSON(path string) (List, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return List{}, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return List{}, err
	}

	var list List
	err = json.Unmarshal(byteValue, &list)
	if err != nil {
		return List{}, err
	}

	return list, nil
}

// Printer functions

func (node Node) Stringify(i int) string {
	s := fmt.Sprintf("%*sNode: %x\n", i, "", hashNode(node))
	s += fmt.Sprintf("%*sPrev: %x\n", i+2, "", node.PrevHash)
	s += fmt.Sprintf("%*sData: %x\n", i+2, "", node.Data)

	return s
}

func (list List) Stringify() string {
	s := fmt.Sprintf("HashedList:\n")
	s += fmt.Sprintf("  Queued Nodes:\n")
	for _, node := range list.queuedNodes {
		s += node.Stringify(4) + "\n"
	}
	s += fmt.Sprintf("  Validated Nodes:\n")
	for _, node := range list.nodes {
		s += node.Stringify(4) + "\n"
	}
	return s
}
