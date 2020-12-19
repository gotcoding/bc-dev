package pos

import (
	"fmt"
	"time"
)

// 使用切片存储链
type Chain struct {
	Blocks []*Block
}

func (c *Chain) Print() {
	for _, v := range c.Blocks {
		fmt.Println(v)
	}
}

func NewChain() *Chain {
	// 创世区块
	block := &Block{
		PreHash:     "0000000000000000000000000000000000000000000000000000000000000000",
		TimeStamp:   time.Now().Format("2006-01-02 15:04:05"),
		Index:       0,
		Validator:   "0000000000",
		Transaction: "创世区块",
	}
	block.CalculateHash()
	c := &Chain{
		Blocks: make([]*Block, 0),
	}
	c.Blocks = append(c.Blocks, block)
	return c
}
