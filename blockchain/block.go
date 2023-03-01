package blockchain

import (
	"bytes"
	"crypto/sha256"
	"lighteningchain/utils"
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte //区块hash值就是其ID
	PrevHash  []byte
	Data      []byte
	Nonce     int64
	Target    []byte
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.Int64ToByte(b.Timestamp),
		b.PrevHash, b.Target, utils.Int64ToByte(b.Nonce), b.Data}, []byte{})
	hash := sha256.Sum256(information) //软件包sha256 实现 FIPS 180-4 中定义的 SHA224 和 SHA256 哈希算法。
	b.Hash = hash[:]
}

func CreateBlock(prevhash []byte, data []byte) *Block {
	block := Block{time.Now().Unix(), []byte{},
		prevhash, []byte{}, 0, data}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash() //所有数据添加好后再计算hash
	return &block
}

func GenesisBlock() *Block {
	genesisWords := "HelloWorld!"
	return CreateBlock([]byte{}, []byte(genesisWords))
}
