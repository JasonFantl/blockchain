package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type Wallet struct {
	publicKey, privateKey string
}

func generateWallet() (Wallet, error) {
	privateKey, publicKey := GenerateRsaKeyPair()

	strPublicKey, err := ExportRsaPublicKeyAsPemStr(publicKey)
	if err != nil {
		return Wallet{}, err
	}
	strPrivateKey := ExportRsaPrivateKeyAsPemStr(privateKey)

	return Wallet{strPublicKey, strPrivateKey}, nil
}

func (wallet Wallet) send() {

}

func (wallet Wallet) signMessage(message []byte) ([]byte, error) {
	privateKey, err := ParseRsaPrivateKeyFromPemStr(wallet.privateKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(message)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, newhash, hashed, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return signature, nil
}

func checkSignature(message, signature []byte, publicKeyString string) bool {
	publicKey, err := ParseRsaPublicKeyFromPemStr(publicKeyString)
	if err != nil {
		return false
	}
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(message)
	hashed := pssh.Sum(nil)

	err = rsa.VerifyPSS(publicKey, newhash, hashed, signature, nil)

	if err != nil {
		return false
	}
	return true
}

func stringifyWallet(w Wallet) string {
	return fmt.Sprintf("Wallet:\n  Public key: %s\n  Private key: %s\n", w.publicKey, w.privateKey)
}
