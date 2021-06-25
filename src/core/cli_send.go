package core

import (
	"fmt"
	"log"
)

//send方法
func (cli *CLI) send(from,to string,amount int,sendname,nodeID string, mineNow bool) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Address is not valid")
	}
	if from ==to {
		log.Panic("ERROR: 发送方和接收方的地址不能相同")
	}//解决余额翻倍漏洞

	bc := NewBlockchain(nodeID)//得到数据库文件存储的区块链并进行实例化,
	UTXOSet := UTXOSet{bc}//创建一个结构体，代表已有区块链的UTXO集
	defer bc.Db().Close()

	tx := NewUTXOTransaction(from, to,sendname,nodeID, amount, &UTXOSet)//创建一笔未花费输出交易,先寻找输出地址(to)为本次发送地址(from)的交易输出,
	// 然后统计余额是否大于amount,然后生成本次交易的输出交易(输出地址为to,余额为amount的交易,以及可能的对于from的交易(统计的余额大于amount))

	if bc.VerifyTransaction(tx) != true {	//在一笔交易被放入一个块之前进行验证
		log.Panic("ERROR: 无效 transaction")
	}

	if mineNow {
		cbTx := NewCoinbaseTX(from, "")//针对from创建一个挖矿输出,奖励50,data是附加数据,输出一个交易transaction
		txs := []*Transaction{cbTx, tx}//一个区块中的多个交易,挖矿输出交易,还有就是转账交易
		newBlock := bc.MineBlock(txs)//针对交易生成一个新的区块,并添加到数据库中
		UTXOSet.Update(newBlock)//当区块链中的区块增加后，要同步更新UTXO集,这里引入的区块为新加入的区块
		fmt.Println("交易发送成功,新生成的交易添加到区块链中.")
	} else {
		sendTx(knownNodes[0], tx)
		fmt.Println("交易发送成功,新生成的交易添加到内存池中.")
	}
}


func (cli *CLI) Send(from,to string,amount int,sendname,nodeID string, mineNow bool) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockchain(nodeID)//得到数据库文件存储的区块链并进行实例化,
	UTXOSet := UTXOSet{bc}//创建一个结构体，代表已有区块链的UTXO集
	defer bc.Db().Close()

	tx := NewUTXOTransaction(from, to,sendname,nodeID, amount, &UTXOSet)//创建一笔未花费输出交易,先寻找输出地址(to)为本次发送地址(from)的交易输出,
	// 然后统计余额是否大于amount,然后生成本次交易的输出交易(输出地址为to,余额为amount的交易,以及可能的对于from的交易(统计的余额大于amount))

	if mineNow {
		cbTx := NewCoinbaseTX(from, "")//针对from创建一个挖矿输出,奖励50,data是附加数据,输出一个交易transaction
		txs := []*Transaction{cbTx, tx}//一个区块中的多个交易,挖矿输出交易,还有就是转账交易
		newBlock := bc.MineBlock(txs)//针对交易生成一个新的区块,并添加到数据库中
		UTXOSet.Update(newBlock)//当区块链中的区块增加后，要同步更新UTXO集,这里引入的区块为新加入的区块
	} else {
		sendTx(knownNodes[0], tx)
	}
	fmt.Println("交易发送成功,新生成的交易添加到内存池中.")
}

