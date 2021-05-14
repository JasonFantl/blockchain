package hashedList

type Node struct {
	Data     []byte
	PrevHash string
	Nonce    int
}

func (node *Node) SetData(data []byte) {
	node.Data = data
}

func (node *Node) AppendData(data []byte) {
	node.Data = append(node.Data, data...)
}

func verifyNode(node Node) bool {
	return verifyHash(hashNode(node))
}

func verifyHash(hash string) bool {
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

func hashNode(node Node) string {
	bytes, err := getBytes(node)
	if err != nil {
		return ""
	}

	return hash(bytes)
}
