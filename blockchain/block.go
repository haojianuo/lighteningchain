package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"lighteningchain/transaction"
	"lighteningchain/utils"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte //区块hash值就是其ID
	PrevHash     []byte
	Nonce        int64
	Target       []byte
	Transactions []*transaction.Transaction
}

func (b *Block) BackTXSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.Int64ToByte(b.Timestamp),
		b.PrevHash, b.Target, utils.Int64ToByte(b.Nonce), b.BackTXSummary()}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

func CreateBlock(prevhash []byte, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{},
		prevhash, 0, []byte{}, txs}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash() //所有数据添加好后再计算hash
	return &block
}

func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte("This is nothing"), []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.Handle(err)
	return res.Bytes()
}

func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
