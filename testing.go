package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jasonfantl/cryptocurrency/blockchain"
)

func testBlockchain() {

	fmt.Println("generating wallet...")
	wallet, err := blockchain.GenerateWallet()
	fmt.Println("generated wallet. Joining network...")

	bc := blockchain.NewBlockchain()
	fmt.Println("Joined network")

	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text == "quit" || text == "exit" {
			return
		} else if text == "TXs" {
			sums := bc.GetSums()
			for name, sum := range sums {
				fmt.Println(string(blockchain.Shorten([]byte(name))), sum)
			}
		} else {
			go bc.Send(wallet, []byte(`-----BEGIN RSA PUBLIC KEY-----
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
j63NkNkp/ZjbFG5Xfa0hUfx/oIkNXoHzoubWPon+4aoVOca659xdXOyvqOAzzJ3J
7W5RlTUmX/H4YYaJkilEGpma7dIktR0UdTdNC5rWnXdn/91x2lWzUKh9gCpcJTld
6yw2OMadItpcRSg9PVZTzbNcU6ardgwdjQpSzYthYVGo3FudA1x0nykIaaLgQkWu
35Uqg/y3cldjeZ/nVzruFmmprbf2U1QY2HY7Ftj49SlhOp5Io57BjbqUjwRep2sH
MXEsXUGn4YSed2ZaqRYvB1p8Djn7kM0fvWvSSSY1DfTpxDMoCw2WYEdYV+SmfOse
CLrf5/wY3goGCIdP6EDY7SP7OiRsQjj5dYuQEwhXHGmAf0l2hpOaPPnSA5XuLiIZ
itMIdA3Ys6XdoIb+DpmmE1JTw1ZjlZBbtP6c/3iMMC+bO2eLiH2wpyYpWVQjEBm/
9PahTCeyyMDtUE3q7ehoXVgYXSEI1nGOQ6FZO9ZHnlddiNVP6VNK3UNYtOh6GJ31
iwb/JHDJuNWCQuoKBGEg8vHs5wY4DmII+VgFrpDloVnW2hrDjt8ChFPbJuhZ2OLL
CIVhIvRIFXgTCs9IM3lSXBMNuHsE2kW7ELuQbymsli+oiZ28/0fYVidyeAw5yWZT
xzVPc8NUUFzqbzFuoR0NGmsCAwEAAQ==
-----END RSA PUBLIC KEY-----`), 30)
		}
	}
}
