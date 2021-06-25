package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/boltdb/bolt"
	"errors"
	"fmt"
	"log"
	"os"
	// Bolt是一个纯粹Key/Value模型的程序。该项目的目标是为不需要完整数据库服务器（如Postgres或MySQL）的项目提供一个简单，快速，可靠的数据库。
)


const dbFile = "blockchain_%s.db"//数据库名
const blocksBucket = "blocks"//表名
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip []byte//区块链里的最后一个区块hash,作为键时,它指向的值为最新的区块,作为值时它的键为"l"
	db *bolt.DB//BoltDB数据库(go语言实现的简介数据库,不需要运行一个服务器,允许我们构建想要的数据结构),采用键值对存储,存储在bucket中.
	//Bolt数据库没有数据类型,键和值都是字节数组.鉴于需要在里面存储GO的结构体,需要对他们进行序列化,也就是说实现GO struct和byte array的互相转换(采用encoding/gob来完成)
}//tip是存储最后一个块的哈希;db存储数据库链接


func(bc *Blockchain) Db() *bolt.DB {
	return bc.db
}

//把区块添加进区块链,挖矿
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var block *Block
	var height int64

	//for _, tx := range transactions {
	//	if bc.VerifyTransaction(tx) != true {	//在一笔交易被放入一个块之前进行验证
	//		log.Panic("ERROR: 无效 transaction")
	//	}
	//}
	//只读的方式浏览数据库，获取当前区块链顶端区块的哈希，为加入下一区块做准备
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))	//通过键"l"拿到区块链顶端区块哈希
		currentBlockBytes := b.Get(lastHash)
		block = DeserializeBlcok(currentBlockBytes)
		height = block.Height+1
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//求出新区块
	newBlock := NewBlock(transactions,lastHash,height)

	//把新区块加入到数据库区块链中
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash,newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"),newBlock.Hash)
		if err !=nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}


//在区块链上找到每一个区块中属于address用户的未花费交易输出,返回未花费输出的交易列表
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	//创建一个map，存储已经花费了的交易输出
	spentTXOs := make(map[string][]int)
	//因为要在链上遍历所有区块，所以要使用到迭代器
	bci := bc.Iterator()

	for {
		block := bci.Next()  //迭代

		//遍历当前区块上的交易
		for _,tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //把交易ID转换成string类型，方便存入map中

			//标签
		Outputs:
			//遍历当前交易中的输出切片，取出交易输出
			for outIdx,out := range tx.Vout {
				//在已经花费了的交易输出map中，如果没有找到对应的交易输出，则表示当前交易的输出未花费
				//反之如下
				if spentTXOs[txID] != nil {
					//存在当前交易的输出中有已经花费的交易输出，
					//则我们遍历map中保存的该交易ID对应的输出的index
					//提示：(这里的已经花费的交易输出index其实就是输入TXInput结构体中的Vout字段)
					for _,spentOutIdx := range spentTXOs[txID] {
						//首先要清楚当前交易输出是一个切片，里面有很多输出，
						//如果map里存储的引用的输出和我们当前遍历到的输出index重合,则表示该输出被引用了
						if spentOutIdx == outIdx {
							continue Outputs  //我们就继续遍历下一轮，找到未被引用的输出
						}
					}
				}
				//到这里是得到此交易输出切片中未被引用的输出

				//这里就要从这些未被引用的输出中筛选出属于该用户address地址的输出
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs,*tx)
				}
			}
			//判断是否为coinbase交易
			if tx.IsCoinbase() == false {
				//如果不是,则遍历当前交易的输入
				for _,in := range tx.Vin {
					//如果当前交易的输入是被该地址address所花费的，就会有对应的该地址的引用输出
					//则在map上记录该输入引用的该地址对应的交易输出
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID],in.Vout)
					}
				}
			}
		}
		//退出for循环的条件就是遍历到的创世区块后
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}


//通过对区块进行迭代，返回区块中的所有未花费交易的交易输出集合
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	//var UTXOs []transaction.TXOutput
	UTXO := make(map[string]TXOutputs)

	//创建一个map，存储已经花费了的交易输出
	spentTXOs := make(map[string][]int)
	//因为要在链上遍历区块，所以要使用到迭代器
	bci := bc.Iterator()

	for {
		block := bci.Next()  //迭代

		//遍历当前区块上的交易
		for _,tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //把交易ID转换成16进制编码(常见的hash值的形式),string类型，方便存入map中,

			//标签
		Outputs:
			//遍历当前交易中的输出切片，取出交易输出
			for outIdx,out := range tx.Vout {
				//在已经花费了的交易输出map中，如果没有找到对应的交易输出，则表示当前交易的输出未花费
				//反之如下
				if spentTXOs[txID] != nil {
					//存在当前交易的输出中有已经花费的交易输出，
					//则我们遍历map中保存的该交易ID对应的输出的index
					//提示：(这里的已经花费的交易输出index其实就是输入TXInput结构体中的Vout字段)
					for _,spentOutIdx := range spentTXOs[txID] {
						//首先要清楚当前交易输出是一个切片，里面有很多输出，
						//如果map里存储的引用的输出和我们当前遍历到的输出index重合,则表示该输出被引用了
						if spentOutIdx == outIdx {//outIdx是输出的int值 也就是钱
							continue Outputs  //我们就继续遍历下一轮，找到未被引用的输出
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs,out)//找到的交易中的输出导入到输出切片中
				UTXO[txID] = outs
			}
			//判断是否为coinbase交易
			if tx.IsCoinbase() == false {
				//如果不是,则遍历当前交易的输入
				for _,in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}
		//退出for循环的条件就是遍历到的创世区块后
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	// //遍历交易集合得到交易，从交易中提取出输出字段Vout,从输出字段中提取出属于address的输出
	// for _,tx := range unspentTransactions {
	// 	for _, out := range tx.Vout {
	// 		if out.IsLockedWithKey(pubKeyHash) {
	// 			UTXOs = append(UTXOs,out)
	// 		}
	// 	}
	// }
	//返回未花费交易输出
	return UTXO
}

//func (bc *Blockchain) AddBlock(transactions []*Transaction) {
//	var lastHash []byte
//	var block  *Block
//	var height int64
//
//	err := bc.db.View(func(tx *bolt.Tx) error {//View是BoltDB事务的另一个类型(只读)
//		b := tx.Bucket([]byte(blocksBucket))//打开bloBucke表，获取对象
//		lastHash = b.Get([]byte("l"))//通过Get获得最后一个区块的哈希用来生成新的哈希
//		currentBlockBytes := b.Get(lastHash)
//		block = DeserializeBlcok(currentBlockBytes)
//	    height = block.Height+1
//		return nil
//	})
//
//	if err != nil {
//		log.Panic(err)
//	}
//
//	newBlock := NewBlock(transactions, lastHash,height)
//
//	err = bc.db.Update(func(tx *bolt.Tx) error {//更新区块，添加新的区块
//		b := tx.Bucket([]byte(blocksBucket))
//		err := b.Put(newBlock.Hash, newBlock.Serialize())
//		if err != nil {
//			log.Panic(err)
//		}
//
//		err = b.Put([]byte("l"),newBlock.Hash)
//		if err != nil {
//			log.Panic(err)
//		}
//
//		bc.tip = newBlock.Hash
//
//		return nil
//	})
//}//区块链中添加新的区块


// AddBlock将块保存到区块链中
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlcok(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.tip = block.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}


//找到可以花费的交易输出,这是基于上面的FindUnspentTransactions 方法
func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte,amount int) (int,map[string][]int) {
	//未花费交易输出map集合
	unspentOutputs := make(map[string][]int)
	//未花费交易
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0	//累加未花费交易输出中的Value值

Work:
	for _,tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx,out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID],outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated,unspentOutputs
}


//实例化一个区块链,默认存储了创世区块 ,接收一个地址为挖矿奖励地址
//打开数据文件dbFile,从中读取实例化区块链,如果数据库文件不存在则创建一个文件,读取文件中的区块链表,如果不存在则退出,输出"不存在区块链...."
func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	//打开一个数据库文件，如果文件不存在则创建该名字的文件
	db,err := bolt.Open(dbFile,0600,nil)//返回一个DB类型的结构体和err
	if err != nil {
		log.Panic(err)
	}

	//读写操作数据库
	err = db.Update(func(tx *bolt.Tx) error{//Update是一个读写事务,因为我们要向数据库中添加创世区块
		b := tx.Bucket([]byte(blocksBucket))//Bucket按名称检索桶,参数是字节切片,若不存在则返回nil,存在则返回Bucket结构体
		//查看名字为blocksBucket的Bucket是否存在,获取存储区块的bucket,
		if b == nil {
			//不存在
			fmt.Println("不存在区块链，需要重新创建一个区块链...")
			os.Exit(1)
			//退出导致当前程序以给定的状态代码退出。
			//通常，代码0表示成功，非0表示错误。
			//程序立即终止;不运行延迟函数。
		}
		//如果存在blocksBucket桶，也就是存在区块链
		//通过键"l"映射出顶端区块的Hash值
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip,db}  //此时Blockchain结构体字段已经变成这样了
	return &bc
}

//创建一个新的区块链数据库,address用来接受挖出创世块的奖励
func CreateBlockchain(address,nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	cbtx := NewCoinbaseTX(address,genesisCoinbaseData)//针对address创建一个挖矿输出,奖励50,Data是附加数据
	genesis := NewGenesisBlock(cbtx)//创建创世区块

	db, err := bolt.Open(dbFile,0600,nil)//打开文件名为dbFile的数据库,若没有就创建
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blocksBucket))//创建一个名为blocksBucket的bucket表
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash,genesis.Serialize())//将创世纪块序列化后,与该块的哈希(作为键值)一起存入Bucket
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"),genesis.Hash)//将最近的hash作为值写入名为b的表中,键为"l"
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip,db}

	return &bc
}



//传入一笔交易,找到它引用的交易,并对它进行签名
func (bc *Blockchain) SignTransaction(tx *Transaction,privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _,vin :=range tx.Vin {
		prevTX,err := bc.FindTransaction(vin.Txid) //找到输入引用的输出所在的交易
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey,prevTXs)
}


//通过交易ID找到一个交易,需要迭代所有的区块
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction,error) {
	bci := bc.Iterator()//用于遍历整个区块链

	for {
		block := bci.Next()

		for _,tx := range block.Transactions {
			if bytes.Compare(tx.ID,ID) == 0 {
				return *tx,nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return Transaction{},errors.New("Transaction is not found")
}

//对一笔交易的输入签名进行验证
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX,err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs) //验证签名
}


//返回最新一个区块的区块高度
func (bc *Blockchain)GetBestheight() int64 {
	var block  *Block
	var height int64
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {//通过View方法获取数据
		b := tx.Bucket([]byte(blocksBucket))//打开bloBucke表，获取对象
		lastHash = b.Get([]byte("l"))//通过Get利用值获取value
		currentBlockBytes := b.Get(lastHash)
		block = DeserializeBlcok(currentBlockBytes)
		height = block.Height
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return height
}

// GetBlockHashes 返回链中所有块的哈希列表
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// GetBlock通过其哈希找到一个块并返回它
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block =  *DeserializeBlcok(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

//新添加的代码


//查询dbFile文件是否存储
func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}