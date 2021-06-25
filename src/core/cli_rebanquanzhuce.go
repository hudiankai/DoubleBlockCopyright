package core

import (
	"blockchaincopyright/src/simhash"
	"fmt"
	"log"
	"strconv"
)

func (cli *CLI) ReBanquanzhuce(Useraddress, name,Title,Abstract,ZhuceWorkfile string,Authorizationmoney,Transactionmoney int,nodeID string, mineNow bool) {
	if !ValidateAddress(Useraddress) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockchain(nodeID)//得到数据库文件存储的区块链并进行实例化,
	UTXOSet := UTXOSet{bc}//创建一个结构体，代表已有区块链的UTXO集
	defer bc.Db().Close()

	Content := ReadFile(ZhuceWorkfile)

	hashcontent := simhash.Simhash(simhash.NewWordFeatureSet([]byte(Content)))
	//fmt.Println(hashcontent,"cli中的")


	b := false
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	for {
		block := bci.Next() //从顶端区块向前面的区块迭代
		for _, tx := range block.Transactions {
			if tx.WorkData.Title == Title  {
				if tx.WorkData.Useraddress == Useraddress {
					b = true
					break
				}else if tx.WorkData.Useraddress != Useraddress {
					log.Panic("对应作品的版权归属已经改变,本用户无法修改作品!")
				}
			}
		}

		if b {
			break
		}

		if len(block.PrevBlockHash) == 0 {
			log.Panic("不存在对应作品!")
			break
		}
	}

	hashcontentstr := strconv.FormatUint(hashcontent,10)
	tx := NewBanquanTransaction(Useraddress, name,Title,Abstract,Content,nodeID,hashcontentstr,Authorizationmoney,Transactionmoney, &UTXOSet)//创建一笔未花费输出交易,先寻找输出地址(to)为本次发送地址(from)的交易输出,
	// 然后统计余额是否大于amount,然后生成本次交易的输出交易(输出地址为to,余额为amount的交易,以及可能的对于from的交易(统计的余额大于amount))

	if bc.VerifyTransaction(tx) != true {	//在一笔交易被放入一个块之前进行验证
		log.Panic("ERROR: 无效 transaction")
	}

	if mineNow {
		txs := []*Transaction{tx}//一个区块中的多个交易,挖矿输出交易,还有就是转账交易
		newBlock := bc.MineBlock(txs)//针对交易生成一个新的区块,并添加到数据库中
		UTXOSet.Update(newBlock)//当区块链中的区块增加后，要同步更新UTXO集,这里引入的区块为新加入的区块
	} else {
		sendTx(knownNodes[0], tx)//已知节点knownNodes
		fmt.Println("作品修改交易发送成功,新生成的包含信息的交易进入交易池.")
	}

}

