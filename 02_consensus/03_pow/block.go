package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// Block 区块结构简化版
type Block struct {
	Index     int    // 区块高度
	TimeStamp string // 时间戳
	Diff      int    // 难度值
	PreHash   string // 前一区块的哈希值
	HashCode  string // 当前区块哈希值
	Nonce     int    // 随机数
	Data      string // 交易数据
}

func NewBlock(data string) *Block {
	return &Block{
		Index:     0,
		TimeStamp: time.Now().String(),
		Diff:      4,
		Nonce:     0,
		Data:      data,
	}
}

// 创建创世区块
func GenerateFirstBlock(data string) *Block {
	block := NewBlock(data)
	block.CalculateHash()
	return block
}

// CalculateHash 计算哈希值
func (b *Block) CalculateHash() {
	hashData := strconv.Itoa(b.Index) + b.TimeStamp + strconv.Itoa(b.Diff) + strconv.Itoa(b.Nonce) + b.Data
	hash := sha256.New()
	hash.Write([]byte(hashData))
	hashBytes := hash.Sum(nil)
	b.HashCode = hex.EncodeToString(hashBytes)
}

// 打包区块，并进行挖矿找到符合难度要求的哈希值
func GenerateBlock(data string, preBlock *Block) *Block {
	block := NewBlock(data)
	block.PreHash = preBlock.HashCode
	block.Index = preBlock.Index + 1
	block.HashCode = pow(block.Diff, block)
	return block
}

// Pow工作量证明
func pow(diff int, block *Block) string {
	for {
		// 计算哈希值
		block.CalculateHash()
		// 验证哈希是否符合难度要求
		hash := block.HashCode
		if strings.HasPrefix(hash, strings.Repeat("0", diff)) {
			return hash //挖矿成功
		} else {
			block.Nonce++ //否则调整随机数
		}
	}
}
