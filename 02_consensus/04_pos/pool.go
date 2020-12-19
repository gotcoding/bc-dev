package pos

import (
	"crypto/rand"
	"log"
	"math/big"
)

// 代表挖矿节点
type Node struct {
	Tokens  int    // 质押代币数量
	Days    int    // 质押时间
	Address string // 节点地址
}

func NewNode(tokens, days int, addr string) *Node {
	return &Node{
		Tokens:  tokens,
		Days:    days,
		Address: addr,
	}
}

//挖矿节点
var mineNodesPool []*Node

//概率节点池
var probabilityNodesPool []*Node

// 随机得出挖矿地址（挖矿概率跟代币数量与币龄有关）
func getMinerAddress() string {
	bInt := big.NewInt(int64(len(probabilityNodesPool)))
	// 得出一个随机数，最大不超过随机节点池的大小
	rInt, err := rand.Int(rand.Reader, bInt)
	if err != nil {
		log.Panic(err)
	}
	return probabilityNodesPool[int(rInt.Int64())].Address
}
