package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/jasonfantl/cryptocurrency/blockchain"
)

func testBlockchain() {

	type Transaction struct {
		From, To string
		Amount   int
	}

	bc := blockchain.NewBlockchain()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "quit" || text == "exit" {
			return
		} else if text == "TXs" {
			peers := make(map[string]int, 0)
			for _, data := range bc.Data() {
				var TX Transaction
				err := json.Unmarshal(data, &TX)
				if err != nil {
					continue
				}
				peers[TX.From] -= TX.Amount
				peers[TX.To] += TX.Amount
			}

			fmt.Println(peers)
		} else {
			newTX := Transaction{"a", "b", rand.Intn(100)}
			bytes, err := json.Marshal(newTX)
			if err != nil {
				continue
			}
			err = bc.Append(bytes)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Sent %d from %s to %s\n", newTX.Amount, newTX.From, newTX.To)
		}

	}
}
