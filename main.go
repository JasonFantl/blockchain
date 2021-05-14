package main

import (
	"fmt"
)

type Transaction struct {
	From, To  string
	Amount    int
	Timestamp int
}

func stringifyTransaction(t Transaction, i int) string {
	return fmt.Sprintf("%*sTransaction: %s --> %-10s %d", i, "", t.From, t.To, t.Amount)
}

func main() {
	testNetwork()
	// 	t1, err := GetBytes(Transaction{
	// 		From:   "a",
	// 		To:     "b",
	// 		Amount: 10,
	// 	})
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	t2, err := GetBytes(Transaction{
	// 		From:   "b",
	// 		To:     "c",
	// 		Amount: 6,
	// 	})
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	bc := blockchain.Blockchain{}

	// 	b1 := blockchain.Block{}
	// 	b1.AppendData(t1)
	// 	b1.AppendData(t2)

	// 	bc.QueueBlock(b1)

	// 	b2 := blockchain.Block{}
	// 	t3, err := GetBytes(Transaction{
	// 		From:   "b",
	// 		To:     "c",
	// 		Amount: 6,
	// 	})
	// 	b2.AppendData(t3)

	// 	bc.QueueBlock(b2)

	// 	bc.Mine()

	// 	fmt.Println(bc.Stringify())
	// 	fmt.Println(blockchain.VerifyBlockchain(bc))
	// }

	// func GetBytes(key interface{}) ([]byte, error) {
	// 	var buf bytes.Buffer
	// 	enc := gob.NewEncoder(&buf)
	// 	err := enc.Encode(key)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return buf.Bytes(), nil
}
