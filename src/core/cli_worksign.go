package core

import (
	"log"
)

func (cli *CLI) Worksign(Useraddress, Title,workfile,signame ,nodeID string, mineNow bool) {
	if !ValidateAddress(Useraddress) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockchain(nodeID)//得到数据库文件存储的区块链并进行实例化,
	UTXOSet := UTXOSet{bc}//创建一个结构体，代表已有区块链的UTXO集
	defer bc.Db().Close()

	Content := ReadFile(workfile)

	tx := NewWorksignTransaction(Useraddress, Title,Content,signame,nodeID, &UTXOSet)//

	if bc.VerifyTransaction(tx) != true {	//在一笔交易被放入一个块之前进行验证
		log.Panic("ERROR: 无效 transaction")
	}

	if mineNow {
		txs := []*Transaction{tx}//一个区块中的多个交易,挖矿输出交易,还有就是转账交易
		newBlock := bc.MineBlock(txs)//针对交易生成一个新的区块,并添加到数据库中
		UTXOSet.Update(newBlock)//当区块链中的区块增加后，要同步更新UTXO集,这里引入的区块为新加入的区块
	} else {
		sendTx(knownNodes[0], tx)//已知节点knownNodes
		//fmt.Println("作品内容私钥加密签名成功,内容签名的交易已经进入交易池中.")
	}
}
