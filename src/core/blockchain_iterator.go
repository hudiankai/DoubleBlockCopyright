package core

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代器的结构体
type BlockchainIterator struct {
	currentHash []byte//保存当前区块的hash值
	db *bolt.DB//区块链数据库链接

}


//区块链迭代器的方法
func (bc *Blockchain) Iterator() *BlockchainIterator  {
	return &BlockchainIterator{bc.tip, bc.db}
}//返回一个区块链迭代器的指针


//返回链中的前一个区块
func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {//View是数据库中的一个只读事务。
		b := tx.Bucket([]byte(blocksBucket))
		encodeBlock := b.Get(i.currentHash)//通过键currenthash得到区块信息(序列化之后的字节切片)
		block = DeserializeBlcok(encodeBlock)//解序列化,获得当前块的内容

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash//前一区块的哈希赋值给迭代器中的"当前哈希",继续往上迭代

	return block
}
