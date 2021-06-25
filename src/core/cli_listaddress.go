package core

import (
	"fmt"
	"log"
)

//列出地址名单,钱包集合中的地址有哪些
func (cli *CLI) listAddresses(nodeID string) {
	//实例化一个钱包的映射,用来读取存储钱包数据文件中的数据
	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		fmt.Println("用户地址:",address)
		fmt.Println("用户姓名:",wallets.Wallets[address].Userinformation.Name,"\n")
	}
}


func (cli *CLI) ListAddresses(nodeID string ) {
	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()
	fmt.Println(len(addresses))
	for _, address := range addresses {
		fmt.Println(address)
	}
}