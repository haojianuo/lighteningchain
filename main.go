package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte //区块hash值就是其ID
	PrevHash  []byte
	Data      []byte
}

type BlockChain struct {
	Blocks []*Block
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{Int64ToByte(b.Timestamp), b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(information) //软件包sha256 实现 FIPS 180-4 中定义的 SHA224 和 SHA256 哈希算法。
	b.Hash = hash[:]
}

func Int64ToByte(num int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))
	return buf
}

func CreateBlock(prevhash []byte, data []byte) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, data}
	block.SetHash() //所有数据添加好后再计算hash
	return &block
}

func GenesisBlock() *Block {
	genesisWords := "HelloWorld!"
	return CreateBlock([]byte{}, []byte(genesisWords))
}

func (bc *BlockChain) AddBlock(data string) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, []byte(data))
	bc.Blocks = append(bc.Blocks, newBlock)
}

// CreateBlockChain 初始化区块链
func CreateBlockChain() *BlockChain {
	myBlockchain := BlockChain{}
	myBlockchain.Blocks = append(myBlockchain.Blocks, GenesisBlock())
	return &myBlockchain
}

func main() {
	blockchain := CreateBlockChain()
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
	}
}
