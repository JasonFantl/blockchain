package blockchain

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jasonfantl/cryptocurrency/blockchain/floodNetwork"
)

func testNetwork() {

	recievePacket := func(bs []byte) {
		message := string(bs)
		fmt.Printf("\n%s\n", message)
	}

	logger := func(s string) {
		fmt.Printf("-- %s\n", s)
	}

	n := floodNetwork.New(recievePacket)
	n.SetLogger(logger)

	port := 1234
	joined := false
	counter := 0
	for !joined && counter < 10 {
		joined = n.Join("127.0.0.1:1234", strconv.Itoa(port))
		port++
		counter++
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "quit" || text == "exit" {
			return
		}
		// n.SendMessage(text)
		n.SendMessage([]byte(text))
	}
}
