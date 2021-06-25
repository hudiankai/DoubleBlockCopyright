package core

import (
	"fmt"
)

func (cli *CLI) createWallet(nodeID string) {
	//实例化一个钱包的映射,用来读取钱包文件中存储的数据,然后创建新的钱包,之后一起存入钱包数据文件
	wallets, _ := NewWallets(nodeID)//返回一个结构体,只包含一个string到wallet的映射,并且加载钱包文件中的钱包信息,
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)//钱包文件中原有的钱包以及新创建的钱包一并写入到钱包文件

	fmt.Println("钱包的公钥(字节形式):",wallets.Wallets[address].PublicKey)
	fmt.Printf("新建的钱包地址: %s\n",address)
}


func (cli *CLI) CreateWalletuser(Name, Role, Idcard,password,nodeID string )  {
	//实例化一个钱包的映射,用来读取钱包文件中存储的数据,然后创建新的钱包,之后一起存入钱包数据文件
	wallets, _ := NewWallets(nodeID)//返回一个结构体,只包含一个string到wallet的映射,并且加载钱包文件中的钱包信息,
	fmt.Println(nodeID)
	address := wallets.CreateWalletuser(Name, Role, Idcard,password )
	if nodeID == "3001" || nodeID == "3003" || nodeID == "3004" {
		wallets.SaveToFile("3003")//钱包文件中原有的钱包以及新创建的钱包一并写入到钱包文件
		wallets.SaveToFile("3001")
		wallets.SaveToFile("3004")
	}else {
		wallets.SaveToFile(nodeID)//钱包文件中原有的钱包以及新创建的钱包一并写入到钱包文件
	}

	fmt.Printf("用户的钱包地址: %s\n",address)
}


func (cli *CLI) CreateWallet(nodeID string)  {
	wallets, _ := NewWallets(nodeID)//返回一个结构体,只包含一个string到wallet的映射
	addresses := wallets.GetAddresses()
	fmt.Println(len(addresses))
	for _, address := range addresses {
		fmt.Println(address)
	}

	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)
	add := wallets.GetAddresses()
	fmt.Println(len(add))
	for _, address := range add {
		fmt.Println(address,wallets.Wallets[address].Userinformation)
	}

	fmt.Println("新创建的钱包:",wallets.Wallets[address])
	fmt.Println("钱包的私钥:",wallets.Wallets[address].PrivateKey)
	fmt.Println("钱包的公钥:",wallets.Wallets[address].PublicKey)
	fmt.Println("钱包对应的用户信息:",wallets.Wallets[address].Userinformation)
	fmt.Printf("Your new address: %s\n",address)
}//main函数单独测试用的函数,首字母大写