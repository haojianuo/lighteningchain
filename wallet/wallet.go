package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"lighteningchain/constcoe"
	"lighteningchain/utils"
	"os"
	"path/filepath"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

// NewKeyPair 椭圆曲线密钥对的生成函数
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

func (w *Wallet) Save() {
	filename := filepath.Join(constcoe.Wallets, string(w.Address())+".wlt")

	privKeyBytes, err := x509.MarshalECPrivateKey(&w.PrivateKey)
	utils.Handle(err)
	privKeyFile, err := os.Create(filename)
	utils.Handle(err)
	err = pem.Encode(privKeyFile, &pem.Block{
		// Type:  "EC PRIVATE KEY",
		Bytes: privKeyBytes,
	})
	utils.Handle(err)
	privKeyFile.Close()
}

func LoadWallet(address string) *Wallet {
	filename := filepath.Join(constcoe.Wallets, address+".wlt")
	if !utils.FileExists(filename) {
		utils.Handle(errors.New("no wallet with such address"))
	}

	privKeyFile, err := os.ReadFile(filename)
	utils.Handle(err)
	pemBlock, _ := pem.Decode(privKeyFile)
	utils.Handle(err)
	privKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	utils.Handle(err)
	publicKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: *privKey,
		PublicKey:  publicKey,
	}
}
