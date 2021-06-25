package core

import (
	"fmt"
	"log"
)

//求账户余额
func (cli *CLI) getBalance(address, nodeID string) {
	if !ValidateAddress(address) {//判断输入的地址是否有效,主要是检查后面的校验位是否正确
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}//创建一个结构体，代表UTXO集
	defer bc.Db().Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4] //这里的4是校验位字节数，这里就不在其他包调过来了

	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	//遍历UTXOs中的交易输出out，得到输出字段out.Value,求出余额
	for _,out := range UTXOs {
		balance += out.Value//目标地址的公钥hash对应的交易输出的所有的Value相加便是该地址的余额
	}

	fmt.Printf("Balance of '%s':%d\n",address,balance)//根据挖矿节点的不同,除以相应的倍数
}




func (cli *CLI) GetBalance(address,nodeID string) {
	if !ValidateAddress(address) {//判断输入的地址是否有效,主要是检查后面的校验位是否正确
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}//创建一个结构体，代表UTXO集
	defer bc.Db().Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4] //这里的4是校验位字节数，这里就不在其他包调过来了

	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	//遍历UTXOs中的交易输出out，得到输出字段out.Value,求出余额
	for _,out := range UTXOs {
		balance += out.Value//目标地址的公钥hash对应的交易输出的所有的Value相加便是该地址的余额
	}

	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}
	_wallet := wallets.GetWallet(address)

	fmt.Printf("Balance of '%s':%d\n",_wallet.Userinformation.Name,balance)
}