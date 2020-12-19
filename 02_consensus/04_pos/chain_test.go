package pos

import (
	"testing"
	"time"
)

func init() {
	//手动添加两个节点
	mineNodesPool = append(mineNodesPool, NewNode(1000, 1, "AAAAAAAAAA"))
	mineNodesPool = append(mineNodesPool, NewNode(100, 3, "BBBBBBBBBB"))
	//初始化随机节点池（挖矿概率与代币数量和币龄有关）
	for _, v := range mineNodesPool {
		for i := 0; i <= v.Tokens*v.Days; i++ {
			probabilityNodesPool = append(probabilityNodesPool, v)
		}
	}
}

func TestPos(t *testing.T) {
	chain := NewChain()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		newBlock := GenerateBlock(chain.Blocks[i], "交易信息", "00000")
		chain.Blocks = append(chain.Blocks, newBlock)
	}
	chain.Print()
}
