package main

import (
	"fmt"
	"lighteningchain/blockchain"
	"lighteningchain/transaction"
)

func main() {
	txPool := make([]*transaction.Transaction, 0)
	var tempTx *transaction.Transaction
	var ok bool
	var property int
	chain := blockchain.CreateBlockChain()
	property, _ = chain.FindUTXOs([]byte("Arno"))
	fmt.Println("Balance of Arno: ", property)

	tempTx, ok = chain.CreateTransaction([]byte("Arno"), []byte("Bravo"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)
	property, _ = chain.FindUTXOs([]byte("Arno"))
	fmt.Println("Balance of Arno: ", property)

	tempTx, ok = chain.CreateTransaction([]byte("Bravo"), []byte("Charlie"), 200) // this transaction is invalid
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("Bravo"), []byte("Charlie"), 50)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("Arno"), []byte("Charlie"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)
	property, _ = chain.FindUTXOs([]byte("Arno"))
	fmt.Println("Balance of Arno: ", property)
	property, _ = chain.FindUTXOs([]byte("Bravo"))
	fmt.Println("Balance of Bravo: ", property)
	property, _ = chain.FindUTXOs([]byte("Charlie"))
	fmt.Println("Balance of Charlie: ", property)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}

	//I want to show the bug at this version.

	tempTx, ok = chain.CreateTransaction([]byte("Bravo"), []byte("Charlie"), 30)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("Bravo"), []byte("Arno"), 30)
	if ok {
		txPool = append(txPool, tempTx)
	}

	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}

	property, _ = chain.FindUTXOs([]byte("Arno"))
	fmt.Println("Balance of Arno: ", property)
	property, _ = chain.FindUTXOs([]byte("Bravo"))
	fmt.Println("Balance of Bravo: ", property)
	property, _ = chain.FindUTXOs([]byte("Charlie"))
	fmt.Println("Balance of Charlie: ", property)
}
