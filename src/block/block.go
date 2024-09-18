package block

import (
	. "gochain/cryptor"
	. "gochain/tx"
	"strconv"
	"time"
)

// // Block 是一个通用区块接口，适用于不同的共识机制
// type Block interface {
// 	GetHash() string         // 返回区块的哈希值
// 	GetPreviousHash() string // 返回前一个区块的哈希值
// 	GetTimestamp() int64     // 返回区块的时间戳
// 	GetTransList() []Tx      // 返回区块中的交易
// 	GetHeight() int          // 长度
// 	Validate() bool          // 验证区块的有效性
// }

// // PoWBlock 定义了工作量证明（PoW）区块的特定接口
// type PoWBlock interface {
// 	Block               // 继承通用区块接口
// 	GetNonce() int      // 返回区块中的 nonce
// 	GetDifficulty() int // 返回区块的挖矿难度
// 	MineBlock()         // 进行挖矿（计算 nonce）
// }

type PoWBlock struct {
	Height       int    `json:"height"`
	PreviousHash string `json:"previousHash"`
	Hash         string `json:"hash"`
	Timestamp    int64  `json:"timestamp"`
	Nonce        int    `json:"nonce"`
	Difficulty   int    `json:"difficulty"`
	MerkleRoot   string `json:"merkleRoot"`
	TransList    []Tx   `json:"transList"`
}

func NewPoWBlock(height int, previousHash string, transList []Tx) *PoWBlock {
	b := &PoWBlock{
		Height:       height,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		TransList:    transList,
	}
	b.MerkleRoot = b.calculateMerkleRoot()
	return b
}

func (b *PoWBlock) MineBlock(difficulty int) {
	b.Difficulty = difficulty
	b.Nonce = 1
	var hash string
	for {
		hash = b.calculateHash()
		if validatePoW(hash, difficulty) {
			break
		}
		b.Nonce++
	}
	b.Hash = hash
}

func (b *PoWBlock) Validate() bool {
	// 验证melkle root
	if b.MerkleRoot != b.calculateMerkleRoot() {
		return false
	}
	// 验证PoW
	if !validatePoW(b.Hash, b.Difficulty) {
		return false
	}
	// 验证hash
	if b.Hash != b.calculateHash() {
		return false
	}
	//验证交易
	for _, tx := range b.TransList {
		if !tx.Validate() {
			return false
		}
	}
	return true
}

func (b *PoWBlock) calculateHash() string {
	data := strconv.Itoa(b.Height) + b.PreviousHash + strconv.FormatInt(b.Timestamp, 10) + strconv.Itoa(b.Nonce) + strconv.Itoa(b.Difficulty) + b.MerkleRoot
	return Sha256(data)
}

func validatePoW(hash string, difficulty int) bool {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}
	return hash[:difficulty] == prefix
}

func txHash(t Tx) string {
	data := t.From + t.To + strconv.Itoa(t.Value) + strconv.FormatInt(t.Timestamp, 10) + t.Signature + t.PubKey + t.Data
	return Sha256(data)
}

/*不是真正的melkle，待优化*/
func (b *PoWBlock) calculateMerkleRoot() string {
	tLen := len(b.TransList)

	if tLen == 0 {
		return ""
	} else if tLen == 1 {
		return txHash(b.TransList[0])
	}

	var merkle string = ""
	for _, tx := range b.TransList {
		merkle = Sha256(merkle + txHash(tx))
	}
	return merkle
}
