package blockchain

import (
	"fmt"
	"strconv"

	"github.com/jasonfantl/cryptocurrency/blockchain/floodNetwork"
)

type Blockchain struct {
	chain   ListChain
	network floodNetwork.Network
}

func NewBlockchain() *Blockchain {
	b := Blockchain{}
	b.chain = ListChain{}
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

func (bc *Blockchain) append(tx Transaction, wallet Wallet) error {
	err := bc.chain.Append(tx, wallet)
	if err != nil {
		return err
	}
	return bc.announceChain()
}

func (bc *Blockchain) Send(wallet Wallet, to []byte, amount int) {
	TX := Transaction{
		From:   wallet.publicKey,
		To:     to,
		Amount: amount,
	}
	fmt.Println("sending " + TX.String())

	err := bc.append(TX, wallet)
	if err != nil {
		fmt.Println(err)
	}
}

func (bc Blockchain) GetSums() map[string]int {
	sums := make(map[string]int)
	for _, nodes := range bc.chain.Nodes {
		for _, node := range nodes {
			sums[string(node.Signed.Transaction.From)] -= node.Signed.Transaction.Amount
			sums[string(node.Signed.Transaction.To)] += node.Signed.Transaction.Amount
		}
	}

	return sums
}
