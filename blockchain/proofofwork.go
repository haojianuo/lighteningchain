package blockchain

import (
	"bytes"
	"crypto/sha256"
	"lighteningchain/constcoe"
	"lighteningchain/utils"
	"math"
	"math/big"
)

func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-constcoe.Difficulty))
	return target.Bytes()
}

func (b *Block) GetDataBaseNonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.Int64ToByte(b.Timestamp),
		b.PrevHash,
		utils.Int64ToByte(nonce),
		b.Target,
		b.BackTXSummary(),
	},
		[]byte{},
	)
	return data
}

func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTarget big.Int

	intTarget.SetBytes(b.Target)

	var hash [32]byte
	var nonce int64
	nonce = 0

	for nonce < math.MaxInt64 {
		data := b.GetDataBaseNonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

func (b *Block) ValidatePoW() bool {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)
	data := b.GetDataBaseNonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	if intHash.Cmp(&intTarget) == -1 {
		return true
	}
	return false
}
