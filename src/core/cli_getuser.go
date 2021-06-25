package core

import (
	"fmt"
	"log"
)

func (cli *CLI) Getuser(username,userpassword,nodeID string) {

	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	for address ,wallet := range wallets.Wallets {
		if wallet.Userinformation.Name == username && wallet.Userinformation.Password == userpassword {
			fmt.Println("用户地址:",address)
			fmt.Println("用户姓名:",wallet.Userinformation.Name,"用户角色:",wallet.Userinformation.Role,"用户IDcard:",wallet.Userinformation.Idcard,"注册时间:",wallet.Userinformation.Time)
		}
	}
}

