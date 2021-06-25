package core

import (
	"crypto/sha256"

)

//创建一个merkle树结构体
type MerkleTree struct {
	RootNode *MerkleNode//根节点
}

//merkle树的节点
type MerkleNode struct {
	Left 	*MerkleNode
	Right 	*MerkleNode
	Data 	[]byte//数据
}

//创建一个新的树节点
func NewMerkleNode(left,right *MerkleNode,data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		//叶子节点
		hash := sha256.Sum256(data)//序列化的交易传入
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data,right.Data...)//当一个节点被关联到其他节点,把其他节点的数据取过来连接
		hash := sha256.Sum256(prevHashes)//连接后再hash
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}

//创建生成一颗新树
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	//输入的交易个数如果是单数的话，就复制最后一个，成为复数
	if len(data) % 2 != 0 {
		data = append(data,data[len(data) - 1])
	}


	for _,datum := range data {
		node := NewMerkleNode(nil,nil,datum)//将数据转化成叶子节点
		nodes = append(nodes,*node)
	}
	
	//循环一层一层的生成节点，知道到最上面的根节点为止
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j],&nodes[j+1],nil)//将叶子节点生成树
			newLevel = append(newLevel,*node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree

}
