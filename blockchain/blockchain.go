package blockchain

type BlockChain struct {
	Blocks []*Block
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
