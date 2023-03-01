package main

import (
	"fmt"
	"lighteningchain/blockchain"
	"time"
)

func main() {
	blockchain := blockchain.CreateBlockChain()
	time.Sleep(time.Second)
	blockchain.AddBlock("This is first Block after Genesis")
	time.Sleep(time.Second)
	blockchain.AddBlock("This is second!")
	time.Sleep(time.Second)
	blockchain.AddBlock("Awesome!")
	time.Sleep(time.Second)

	for num, block := range blockchain.Blocks {
		fmt.Printf("number:%d Timestamp: %d\n", num, block.Timestamp)
		fmt.Printf("number:%d hash: %x\n", num, block.Hash)
		fmt.Printf("number:%d Previous hash: %x\n", num, block.PrevHash)
		fmt.Printf("number:%d data: %s\n", num, block.Data)
		fmt.Printf("number:%d nonce:%d\n", num, block.Nonce)
		fmt.Println("POW validation:", block.ValidatePoW())

	}
}
