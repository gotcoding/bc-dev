package pos

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Block struct {
	PreHash     string // 前一区块哈希
	HashCode    string // 区块哈希
	TimeStamp   string // 时间戳
	Index       int    // 区块高度
	Validator   string // 验证者
	Transaction string // 交易数据
}

// CalculateHash 计算哈希值
func (b *Block) CalculateHash() {
	hashData := strconv.Itoa(b.Index) + b.PreHash + b.TimeStamp + b.Validator + b.Transaction
	hash := sha256.New()
	hash.Write([]byte(hashData))
	hashBytes := hash.Sum(nil)
	b.HashCode = hex.EncodeToString(hashBytes)
}

// 打包区块
func GenerateBlock(preBlock *Block, trans string, addrs string) *Block {
	block := &Block{
		PreHash:     preBlock.HashCode,
		TimeStamp:   time.Now().Format("2006-01-02 15:04:05"),
		Index:       preBlock.Index + 1,
		Validator:   addrs,
		Transaction: trans,
	}
	block.Validator = getMinerAddress()
	block.CalculateHash()
	return block
}
