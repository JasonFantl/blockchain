package blockchain

import (
	"strconv"

	"github.com/jasonfantl/cryptocurrency/blockchain/chainStruct"
	"github.com/jasonfantl/cryptocurrency/blockchain/floodNetwork"
)

type Blockchain struct {
	chain   chainStruct.Chain
	network floodNetwork.Network
}

func NewBlockchain() *Blockchain {
	b := Blockchain{}
	b.chain = new(chainStruct.SimpleChain)
	b.network = floodNetwork.New(b.recievePacket)

	// join network. First find open port to init from
	port := 1234
	joined := false
	counter := 0
	for !joined && counter < 10 {
		joined = b.network.Join("127.0.0.1:1234", strconv.Itoa(port))
		port++
		counter++
	}

	return &b
}

func (bc *Blockchain) Append(data []byte) error {
	// return bc.announceBlock(bc.list.Append(data))
	bc.chain.Append(data)
	return bc.announceChain()
}

func (bc *Blockchain) Data() [][]byte {
	return bc.chain.Array()
}
