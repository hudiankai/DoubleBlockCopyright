package core

import (
	"fmt"
	"log"
)

func (cli *CLI) Shouquan(Buyaddress,buyname,Selladdress, ShouquanTitle,nodeID string, mineNow bool) {
	//验证两个地址
	if !ValidateAddress(Buyaddress) {
		log.Panic("ERROR: BuyAddress is not valid")
	}
	if !ValidateAddress(Selladdress) {
		log.Panic("ERROR: SellAddress is not valid")
	}

	bc := NewBlockchain(nodeID)//得到数据库文件存储的区块链并进行实例化,
	UTXOSet := UTXOSet{bc}//创建一个结构体，代表已有区块链的UTXO集
	defer bc.Db().Close()

	b := false
	var workdata Work
	var blockhash []byte
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	for {

		block := bci.Next()	//从顶端区块向前面的区块迭代

		for _,tx := range block.Transactions {
			if tx.WorkData.Title == ShouquanTitle {
				if Selladdress == tx.WorkData.Useraddress {
					blockhash = block.Hash
					workdata = tx.WorkData
					b = true
					break
				}else {
					log.Panic("作品版权所有权已经更改,请向最新的版权所有者请求版权交易.")
				}
			}
		}

		if b {
			break
		}//找到对应的作品并且版权归属正确

		if len(block.PrevBlockHash) == 0 || b == true {//遍历完全部区块或者是找到对应的作品
			if b != true {
				workdata = Worknil
			}
			break
		}
	}


	tx := NewShouquanTransaction(Buyaddress,buyname,Selladdress, ShouquanTitle,nodeID, &UTXOSet,workdata,blockhash)

	if bc.VerifyTransaction(tx) != true {	//在一笔交易被放入一个块之前进行验证
		log.Panic("ERROR: 无效 transaction")
	}
	//fmt.Println("clishouquan中的63行,刚生成tx",tx.ShouquanData.txhash,tx.ShouquanData.blockhash)

	if mineNow {
		txs := []*Transaction{tx}//一个区块中的多个交易,挖矿输出交易,还有就是转账交易
		newBlock := bc.MineBlock(txs)//针对交易生成一个新的区块,并添加到数据库中
		UTXOSet.Update(newBlock)//当区块链中的区块增加后，要同步更新UTXO集,这里引入的区块为新加入的区块
	} else {
		sendTx(knownNodes[0], tx)//已知节点knownNodes
		fmt.Println("作品使用权购买交易发送成功,新生成的包含该信息的交易添加到交易池中.")
	}
	//版权注册成功之后,在work.date文件中生成work存储信息
}
