package blockchain

import (
	"bytes"
	"encoding/gob"
	"lighteningchain/constcoe"
	"lighteningchain/transaction"
	"lighteningchain/utils"
	"os"
)

type TransactionPool struct {
	Txs []*transaction.Transaction
}

func (tp *TransactionPool) AddTransaction(tx *transaction.Transaction) {
	tp.Txs = append(tp.Txs, tx)
}

func (tp *TransactionPool) SaveFile() {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(tp)
	utils.Handle(err)
	err = os.WriteFile(constcoe.TransactionPoolFile, content.Bytes(), 0644)
	utils.Handle(err)
}

func (tp *TransactionPool) LoadFile() error {
	if !utils.FileExists(constcoe.TransactionPoolFile) {
		return nil
	}

	var transactionPool TransactionPool

	fileContent, err := os.ReadFile(constcoe.TransactionPoolFile)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&transactionPool)

	if err != nil {
		return err
	}

	tp.Txs = transactionPool.Txs
	return nil
}

func CreateTransactionPool() *TransactionPool {
	transactionPool := TransactionPool{}
	err := transactionPool.LoadFile()
	utils.Handle(err)
	return &transactionPool
}

func RemoveTransactionPoolFile() error {
	err := os.Remove(constcoe.TransactionPoolFile)
	return err
}
