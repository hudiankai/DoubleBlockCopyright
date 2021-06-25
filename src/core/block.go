package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob" //数据结构序列化的编码/解码工具。编码使用Encoder，解码使用Decoder。
	"log"
	"time"
)

type Block struct {
	Timestamp int64//区块创建的时间
	Transactions []*Transaction//区块交易
	PrevBlockHash []byte//前一个区块的hash值
	Hash []byte//区块当前的hash值，用于校验数据有效
	Nonce int//工作量难度,计数器(密码学术语),对工作量证明进行验证时用到
	Height int64//区块的高度
}//结构体的属性必须大写,不然没法序列化


//序列化： 将数据结构或对象转换成二进制串的过程。
//反序列化：将在序列化过程中所生成的二进制串转换成数据结构或者对象的过程。
//将block序列化为一个字节数组,
func (b *Block) Serialize() []byte  {
	var result bytes.Buffer//定义字节buffer，存储序列化后的区块数据.//Buffer缓冲区是具有读和写方法的可变大小的字节缓冲区。
	//缓冲区的零值是一个可以使用的空缓冲区。
	encoder := gob.NewEncoder(&result)//创建基于buffer内存的编码器

	err := encoder.Encode(b)//序列化,使用编码器对block结构体进行编码
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()//返回result的字节数组
}//将区块链序列化成字节切片

func NewBlock(transactions []*Transaction,preBlockHash []byte,height int64) *Block{
	block := &Block{time.Now().Unix(),transactions,preBlockHash,[]byte{},0,height}
    pow :=NewProofOfWork(block)
    nonce,hash := pow.Run()//调用计算哈希的方法,返回计数器和区块的hash值

    block.Hash = hash[:]
    block.Nonce = nonce

    //fmt.Println("挖矿成功的时候的区块以及pow验证:")
    //fmt.Println("hashtran如下:",block.HashTransactions(),"\n",block.PrevBlockHash,"\n","交易的个数",len(block.Transactions))
	//boolen := pow.Validate()
	//fmt.Printf("POW is %s\n",strconv.FormatBool(boolen))
	//fmt.Println("验证成功是区块中的交易")
	//for _,tx :=range block.Transactions {
	//	fmt.Println(tx)
	//}


	return block
}//通过交易transaction切片和前一区块的hash值创建新的区块并返回,


//merkle树的形式拼接字符串
//func (b *Block) HashTransactions() []byte {
//
//	var transactions  [][]byte
//
//	for _,tx := range b.Transactions {
//		transactions = append(transactions,tx.Serialize())//先获得每笔交易的哈希在连接起来
//	}
//	mTree := NewMerkleTree(transactions)//使用序列化后的交易构建一个merkle树
//
//	return mTree.RootNode.Data//树根作为块交易的唯一标识符
//}

//交易的字符串拼接验证
func (b *Block) HashTransactions() []byte {
	var txHash [32]byte
	var txHashes [][]byte
	for _,tx := range b.Transactions {
		txHashes = append(txHashes,tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))

	return txHash[:]
}

//区块内部交易的存储验证
func (b *Block) MineHashTransactions() []byte {

	var transactions  bytes.Buffer

	for _,tx := range b.Transactions {
		transactions.Write(tx.Serialize())//先获得每笔交易的哈希在连接起来
	}

	hash := sha256.Sum256(transactions.Bytes())

	return hash[:] //树根作为块交易的唯一标识符
}

func NewGenesisBlock(coinbase *Transaction) *Block  {

	return NewBlock([]*Transaction{coinbase},[]byte{},1)//coinbase是一个交易transaction

}//创建创世区块,创世区块只能创建一次


//将字节反序列化为一个block结构体
func DeserializeBlcok(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))//创建解码器,初始化反序列化对象

	err := decoder.Decode(&block)//通过Decode（）进行反序列化,对于d内容进行解码,解码的内容写入变量block的内存中
	if err != nil {
		log.Panic(err)
	}

	return &block
}//反序列化一个区块，传入字节数组，返回区块

