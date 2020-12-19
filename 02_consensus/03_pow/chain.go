package main

import "fmt"

// 区块节点
type Node struct {
	NextNode *Node
	Data     *Block
}

// 采用链表形式维护区块链
type Chain struct {
	Head *Node
	Tail *Node
	Len  int
}

// 新建一条链
func NewChain(firstBlock *Block) *Chain {
	node := NewNode(firstBlock)
	return &Chain{
		Head: node,
		Tail: node,
		Len:  1,
	}
}

func NewNode(b *Block) *Node {
	return &Node{
		NextNode: nil,
		Data:     b,
	}
}

func (c *Chain) AddNode(b *Block) {
	node := NewNode(b)
	c.Len++
	c.Tail.NextNode = node
	c.Tail = node
}

// 遍历链表，打印区块链
func (c *Chain) Print() {
	if c.Len == 0 {
		fmt.Println("Empty Chain")
	}
	cur := c.Head
	for cur != nil {
		fmt.Printf("index: %d, preHash: %s, hash: %s, nonce: %d, data: %s \n", cur.Data.Index,
			cur.Data.PreHash, cur.Data.HashCode, cur.Data.Nonce, cur.Data.Data)
		cur = cur.NextNode
	}
}
