package block

import (
	. "gochain/tx"
)

// Block 是一个通用区块接口，适用于不同的共识机制
type Block interface {
	GetHash() string         // 返回区块的哈希值
	GetPreviousHash() string // 返回前一个区块的哈希值
	GetTimestamp() int64     // 返回区块的时间戳
	GetTransList() []Tx      // 返回区块中的交易
	GetHeight() int          // 长度
	Validate() bool          // 验证区块的有效性
}

// PoWBlock 定义了工作量证明（PoW）区块的特定接口
type PoWBlock interface {
	Block               // 继承通用区块接口
	GetNonce() int      // 返回区块中的 nonce
	GetDifficulty() int // 返回区块的挖矿难度
	MineBlock()         // 进行挖矿（计算 nonce）
}

type PoWBlockHeader struct {
	Height       int
	Hash         string
	PreviousHash string
	Timestamp    int64
	Nonce        string
	MerkleRoot   string
}

type SimplePoWBlock struct {
	Height       int
	Hash         string
	PreviousHash string
	Timestamp    int64
	TransList    []Tx
	Nonce        int
	Difficulty   int
}

func (b *SimplePoWBlock) GetHash() string {
	return b.Hash
}

func (b *SimplePoWBlock) GetPreviousHash() string {
	return b.PreviousHash
}

func (b *SimplePoWBlock) GetTimestamp() int64 {
	return b.Timestamp
}

func (b *SimplePoWBlock) GetTransList() []Tx {
	return b.TransList
}

func (b *SimplePoWBlock) Validate() bool {
	// 验证工作量证明
	// return ValidatePoW(b.Hash, b.Difficulty)
	return true
}

func (b *SimplePoWBlock) GetNonce() int {
	return b.Nonce
}

func (b *SimplePoWBlock) GetDifficulty() int {
	return b.Difficulty
}

func (b *SimplePoWBlock) MineBlock() {
	// // 挖矿算法（找到符合难度要求的 nonce）
	// for !IsValidHash(b.Hash, b.Difficulty) {
	// 	b.Nonce++
	// 	b.Hash = CalculateHash(b.PreviousHash, b.Timestamp, b.Transactions, b.Nonce)
	// }
}
