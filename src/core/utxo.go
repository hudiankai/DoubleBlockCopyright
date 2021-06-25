package core

//UTXO（Unspent Transaction Outputs）是未花费的交易输出(输出二字代表某个地址的余额记录,余额为零的则取消记录)，
// 它是比特币交易生成及验证的一个核心概念。 交易构成了一组链式结构，
// 所有合法的比特币交易都可以追溯到前向一个或多个交易的输出(每一笔交易的输入跟输出的金额是一样的,
// 交易输出包含支付者剩余的钱的记录的输出)，这些链条的源头都是挖矿奖励，末尾则是当前未花费的交易输出。
// 账户余额，实际上是钱包通过扫描区块链并聚合所有属于该用户的UTXO(交易的to地址是该账户)计算得来的。
//是从所有区块链交易中构建(只对区块迭代一次)而来的缓存,然后用来计算余额和验证新的交易

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)
const utxoBucket = "chainstate"

//UTXOSet结构表示UTXO集
type UTXOSet struct {
	Blockchain *Blockchain
}

//构建UTXO集的索引并存储在数据库的bucket表中
func (u UTXOSet) Reindex() {
	//调用区块链中的数据库,
	db := u.Blockchain.Db()
	//桶名
	bucketName := []byte(utxoBucket)

	//对数据库进行读写操作
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName) //因为我们是要重新建一个bucket，所以如果原来的数据库中有相同名字的桶，则删除
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}
		//创建新桶
		_,err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//返回链上所有未花费交易中的交易输出
	UTXO := u.Blockchain.FindUTXO()

	//把未花费交易中的交易输出集合写入桶中
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		//写入键值对
		for txID,outs := range UTXO {
			key,err := hex.DecodeString(txID)//字符串转换成字节切片
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key,outs.Serialize())//最终将输出保存在bucket中
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
} 

//查询并返回被用于这次花费的未花费的输出，找到的输出的总额要刚好大于要花费的输入额,
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte,amount int) (int,map[string][]int) {
	//存储找到的未花费输出集合
	unspentOutputs := make(map[string][]int)
	//记录找到的未花费输出中累加的值
	accumulated := 0
	db := u.Blockchain.Db()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		//声明一个游标，类似于我们之前构造的迭代器
		c := b.Cursor()

		//用游标来遍历这个桶里的数据,这个桶里装的是链上所有的未花费输出集合
		for k,v := c.First(); k != nil; k,v =c.Next() {
			txID := hex.EncodeToString(k)//字节切片转成字符串(常见的hash值的样子)
			outs := DeserializeOutputs(v)//反序列化,由字节切片到交易输出的集合

			for outIdx,out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHash) && accumulated < amount {//判断参数公钥是否是验证的交易的公钥
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID],outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated,unspentOutputs
}

//查询对应的地址的未花费输出
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.Db()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))//按照名称检索Bucket,得到一个bucket结构体,Bucket表示数据库中键/值对的集合
		c := b.Cursor()// Cursor游标创建与bucker桶关联的cursor光标,cursor 游标表示一个迭代器，它可以按排序顺序遍历bucket中的所有键/值对。

		for k,v := c.First();k != nil;k,v = c.Next() {//首先将光标移动到bucket中的第一项，并返回它的键和值。 如果桶是空的，那么返回nil键和值。
			outs := DeserializeOutputs(v)//反序列化,由字节切片到交易输出的集合

			for _,out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {//交易输出中存储的公钥hash与检索的目标公钥hash比较
					UTXOs = append(UTXOs,out)//结构体导入结构体切片
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return UTXOs
}

//从区块链中更新 UTXO 集
//当挖出一个新块时，应该更新 UTXO 集。
// **更新意味着移除已花费输出，并从新挖出来的交易中加入未花费输出。**
// 如果一笔交易的输出被移除，并且不再包含任何输出，那么这笔交易也应该被移除。
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.Db()

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _,tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _,vin := range tx.Vin {
					//实例化结构体TXOutputs
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeOutputs(outsBytes)

					for outIdx,out := range outs.Outputs {
						if outIdx != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs,out)
						}
					}
					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)//一个交易中不包含未花费的输出的时候删除它
						if err != nil  {
							log.Panic(err)
						}
					}else{
						err := b.Put(vin.Txid,updatedOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}

			newOutputs := TXOutputs{}
			for _,out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs,out)
			}

			err := b.Put(tx.ID,newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}


//返回UTXO集中的交易数
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.Db() 
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k,_ := c.First(); k != nil; k,_ = c.Next() {
			counter++
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return counter
}


