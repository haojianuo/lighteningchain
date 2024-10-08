package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"lighteningchain/constcoe"
	"lighteningchain/utils"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) CalculateTXHash() []byte {
	var buffer bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(buffer.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID() {
	tx.ID = tx.CalculateTXHash()
}

func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}, nil}
	txOut := TxOutput{constcoe.InitCoin, toaddress}
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

// PlainCopy 创建一个交易的副本，但是不包含签名和公钥
func (tx *Transaction) PlainCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.OutIdx, nil, nil})
	}

	for _, txout := range tx.Outputs {
		outputs = append(outputs, TxOutput{txout.Value, txout.HashPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

func (tx *Transaction) PlainHash(inidx int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevPubKey
	return txCopy.CalculateTXHash()
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsBase() {
		return
	}
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey) // This is because we want to sign the inputs seperately!
		signature := utils.Sign(plainhash, privKey)
		tx.Inputs[idx].Sig = signature
	}
}

func (tx *Transaction) Verify() bool {
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		if !utils.Verify(plainhash, input.PubKey, input.Sig) {
			return false
		}
	}
	return true
}
