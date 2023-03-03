package blockchain

import (
	"encoding/hex"
	"fmt"
	"lighteningchain/transaction"
	"lighteningchain/utils"
)

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// CreateBlockChain 初始化区块链
func CreateBlockChain() *BlockChain {
	myBlockchain := BlockChain{}
	myBlockchain.Blocks = append(myBlockchain.Blocks, GenesisBlock())
	return &myBlockchain
}

func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int)
	//unSpentTxs就是我们要返回包含指定地址的可用交易信息的切片
	//spentTxs用于记录遍历区块链时那些已经被使用的交易信息的Output
	//key值为交易信息的ID值（需要转成string）
	//value值为Output在该交易信息中的序号

	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		block := bc.Blocks[idx]
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			//从最后一个区块开始向前遍历区块链，然后遍历每一个区块中的交易信息
		IterOutputs:
			for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}
				//遍历交易信息的Output，如果该Output在spentTxs中就跳过，说明该Output已被消费。

				if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}
			//否则确认ToAddress正确与否，正确就是我们要找的可用交易信息。
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
			//检查当前交易信息是否为Base Transaction（主要是它没有input）
			//如果不是就检查当前交易信息的input中是否包含目标地址
			//有的话就将指向的Output信息加入到spentTxs中
		}
	}
	return unSpentTxs
}

func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}
	return accumulated, unspentOuts
}
func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}
	return accumulated, unspentOuts
}

func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}
	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{txID, outidx, from}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{acc - amount, from})
	}
	tx := transaction.Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx, true
}

func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}
