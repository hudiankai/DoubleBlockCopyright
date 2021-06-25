package core

import (
	"fmt"
	"log"
)

//创建一条链
func (cli *CLI) createBlockchain(address, nodeID string) {
	if !ValidateAddress(address) {//验证地址是否有效,主要检验后面的校验位
		log.Panic("ERROR: Address is not valid")
	}

	bc := CreateBlockchain(address, nodeID)
	defer bc.Db().Close()

	UTXOSet := UTXOSet{bc}//创建一个UTXOSet结构体
	UTXOSet.Reindex()////构建UTXO集的索引并存储在数据库的bucket表中,会进行UTXO的重新创建,
	fmt.Println("区块链创建成功!")
}


func (cli *CLI) CreateBlockchain(address,nodeID string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := CreateBlockchain(address,nodeID)
	defer bc.Db().Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()
	fmt.Println("Done!")
}