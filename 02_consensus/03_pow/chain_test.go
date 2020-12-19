package main

import (
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	// 构建创世区块和链
	firstBlock := GenerateFirstBlock("first block")
	chain := NewChain(firstBlock)

	// 继续挖几个块
	for i := 1; i < 10; i++ {
		block := GenerateBlock(strconv.Itoa(i), chain.Tail.Data)
		chain.AddNode(block)
	}
	chain.Print()
}
