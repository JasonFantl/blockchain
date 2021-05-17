package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type Transaction struct {
	From, To []byte
	Amount   int
}

func (tx Transaction) String() string {
	return fmt.Sprintf("TX: From %s, To: %s, Amount: %d", Shorten(tx.From), Shorten(tx.To), tx.Amount)
}

func Shorten(b []byte) []byte {
	if len(b) < 105 {
		return b
	}
	return b[100:105]
}

type Wallet struct {
	publicKey, privateKey []byte
}

func (bc Blockchain) GenerateWallet() (Wallet, error) {
	privateKey, publicKey := generateRsaKeyPair()

	strPublicKey, err := exportRsaPublicKeyAsPemStr(publicKey)
	if err != nil {
		return Wallet{}, err
	}
	strPrivateKey := exportRsaPrivateKeyAsPemStr(privateKey)

	return Wallet{strPublicKey, strPrivateKey}, nil
}

func (wallet Wallet) send() {

}

func (wallet Wallet) signMessage(node Node) ([]byte, error) {
	txBytes := hash(node)

	privateKey, err := parseRsaPrivateKeyFromPemStr(wallet.privateKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(txBytes)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, newhash, hashed, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return signature, nil
}

func checkSignature(signed Node, signature, publicKeyBytes []byte) bool {
	txBytes := hash(signed)

	publicKey, err := parseRsaPublicKeyFromPemStr(publicKeyBytes)
	if err != nil {
		return false
	}
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(txBytes)
	hashed := pssh.Sum(nil)

	err = rsa.VerifyPSS(publicKey, newhash, hashed, signature, nil)

	if err != nil {
		return false
	}
	return true
}

func (w Wallet) String() string {
	return fmt.Sprintf("Wallet:\n  Public key: %s\n  Private key: %s\n", w.publicKey, w.privateKey)
}
