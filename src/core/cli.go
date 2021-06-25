package core

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {}


func (cli *CLI) printUsage()  {//打印使用方法
	fmt.Println("Usage:")
	fmt.Println("createwallet -name 用户姓名 -role 用户角色 -idcard 用户身份证号 -password 密码 -Generates a new key-pair and saves it into the wallet file")//创建一个新的公钥私钥对，把它存到文件中
	fmt.Println("listaddress - Lists all addresses from the wallet file")
	fmt.Println("getbalance -name 用户姓名 -password 密码 - Get balance of ADDRESS")
	fmt.Println("createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS ")
	fmt.Println("printchain - print all the blocks  of the blockchain")
	fmt.Println("printlastchain - 打印最新的区块 ")
	fmt.Println("getaddress -getaddressname 用户名 -password 密码  -通过用户名和密码找到地址 ")
	fmt.Println("send -name 姓名 -to TO -amount AMOUNT -password 密码 -mine - Send AMOUNT of coins from FROM address to TO")
	fmt.Println("reindexutxo - Rebuilds the UTXO set")
	fmt.Println("作品版权签名 -name 姓名 -Userpassword 作者密码 -Title 作品名称 -workfile 作品文件名 -mine - 作品内容私钥签名信息")
	fmt.Println("版权注册 -name 姓名 -Userpassword 作者密码 -Title 作品名称 -Abstract 作品摘要 -workfile 作品文件名 " +
		"-Authorizationmoney 作品版权授权费用 -Transactionmoney 作品版权转权费用 -mine - 作品注册版权信息")
	fmt.Println("作品修改 -name 姓名 -Userpassword 作者密码 -Title 作品名称 -Abstract 作品摘要 -workfile 作品文件名 " +
		"-Authorizationmoney 作品版权授权费用 -Transactionmoney 作品版权转权费用 -mine - 可以修改作品的内容以及定价,作品名称无法修改")

	fmt.Println("版权使用权购买 -Selladdress  卖方公钥地址 -Buyname 买方姓名 -buypassword 密码 -Title 作品名称  -mine -作品使用权购买信息")
	fmt.Println("作品版权转让 -Selladdress  版权所有者公钥地址 -Buyname 版权购买者姓名 -buypassword 密码 -Title 作品名称  " +
		"-NewAuthorizationmoney  新的作品版权授权费用 -NewTransactionmoney 新的作品版权转权费用  -mine -作品使用权购买信息")
	fmt.Println("启动节点 -miner ADDRESS - 使用NODE_ID env.var中指定的ID启动节点. -miner enables mining")
	fmt.Println("gettransaction -tranname 用户名 -password 密码 -title 作品题目 -type 类型(认证记录/签名记录/授权记录/转权记录) - 查看作品有关的交易信息(包含作品信息)")
	fmt.Println("getallwork - 展示所有的作品")
	fmt.Println("getworkrecord -title  作品名称 -展示指定作品的所有版权交易信息")
	fmt.Println("getuser -name 查询者姓名 -password 密码 - 查询得到用户信息")
	fmt.Println("createworkfile -filename 文件名 - 创建一个文件")
	fmt.Println("tracebackzhuanquan -title 作品名称 -name 姓名 -password 密码 -版权转让记录溯源")
	fmt.Println("tracebackshouquan -title 作品名称 -name 姓名 -password 密码 -版权使用权授权记录溯源")

}

func (cli *CLI) validateArgs()  {//解析参数
	//os.Args获取程序运行时给出的参数
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}


//负责解析命令行参数和处理命令
func (cli *CLI) Run() {
	//解析命令行额参数并执行命令
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")

		if nodeID == "" {
			fmt.Printf("NODE_ID env. var is not set!")
			os.Exit(1)
		}


	//使用标准数据库里面的flag包来解析命令行参数
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	printlastChainCmd := flag.NewFlagSet("printlastchain",flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	banquanzhuceCmd := flag.NewFlagSet("版权注册", flag.ExitOnError)
	rebanquanzhuceCmd := flag.NewFlagSet("作品修改", flag.ExitOnError)
	worksignCmd := flag.NewFlagSet("作品版权签名", flag.ExitOnError)
	workshouquanCmd := flag.NewFlagSet("版权使用权购买", flag.ExitOnError)
	workzhuanquanCmd := flag.NewFlagSet("作品版权转让", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("启动节点", flag.ExitOnError)
	getaddressCmd :=flag.NewFlagSet("getaddress", flag.ExitOnError)
	gettransactionCmd :=flag.NewFlagSet("gettransaction", flag.ExitOnError)
	getallworkCmd :=flag.NewFlagSet("getallwork", flag.ExitOnError)
	getworkrecordCmd :=flag.NewFlagSet("getworkrecord", flag.ExitOnError)
	getuserCmd := flag.NewFlagSet("getuser",flag.ExitOnError)
	createworkfileCmd := flag.NewFlagSet("createworkfile",flag.ExitOnError)
	tracebackzhuanquanCmd := flag.NewFlagSet("tracebackzhuanquan",flag.ExitOnError)
	tracebackshouquanCmd := flag.NewFlagSet("tracebackshouquan",flag.ExitOnError)



	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")//新添加的mine值
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")

	getBalancename := getBalanceCmd.String("name","","The name to get balacne")
	getBalancepassword := getBalanceCmd.String("password","","The password to get balacne")

	createBlockchainAddress := createBlockchainCmd.String("address","","The address to send genesis block reward to")

	//转账交易的参数
	sendname := sendCmd.String("name","","Source wallet name")
	sendTo := sendCmd.String("to","","Destination wallet address")
	sendAmount := sendCmd.Int("amount",0,"Amount to send")
	sendpassword := sendCmd.String("password","","Source wallet password")

	//作品版权注册所需的参数
	Userpassword := banquanzhuceCmd.String("Userpassword","","the address of auther")
	Title := banquanzhuceCmd.String("Title","","the title of work")
	Abstract := banquanzhuceCmd.String("Abstract","","the Abstract of work")
	Authorizationmoney := banquanzhuceCmd.Int("Authorizationmoney",0,"the Authorizationmoney of work")
	Transactionmoney := banquanzhuceCmd.Int("Transactionmoney",0,"the Transactionmoney of work")
	ZhuceWorkfile := banquanzhuceCmd.String("workfile","","the content of work")
	ZhuceMine := banquanzhuceCmd.Bool("mine", false, "Mine immediately on the same node")//新添加的mine值
	Zhucename :=  banquanzhuceCmd.String("name","","the name of auther")

	//作品版权修改(作品内容和作品定价)
	ReUserpassword := rebanquanzhuceCmd.String("Userpassword","","the address of auther")
	ReTitle := rebanquanzhuceCmd.String("Title","","the title of work")
	ReAbstract := rebanquanzhuceCmd.String("Abstract","","the Abstract of work")
	ReAuthorizationmoney := rebanquanzhuceCmd.Int("Authorizationmoney",0,"the Authorizationmoney of work")
	ReTransactionmoney := rebanquanzhuceCmd.Int("Transactionmoney",0,"the Transactionmoney of work")
	ReZhuceWorkfile := rebanquanzhuceCmd.String("workfile","","the content of work")
	ReZhuceMine := rebanquanzhuceCmd.Bool("mine", false, "Mine immediately on the same node")//新添加的mine值
	ReZhucename :=  rebanquanzhuceCmd.String("name","","the name of auther")

	//用户注册钱包是需要的参数
	WalletName := createWalletCmd.String("name","","the name of wallet")
	WalletRole := createWalletCmd.String("role","","the role of wallet")
	WalletIdcard := createWalletCmd.String("idcard","","the Idcard of wallet")
	WalletPassword := createWalletCmd.String("password","","the password of wallet")

	//用户私钥签名未完成作品内容存证需要的参数
	SignUserpassword := worksignCmd.String("Userpassword","","the address of auther")
	SignTitle := worksignCmd.String("Title","","the title of work")
	Signworkfile := worksignCmd.String("workfile","","the workfile of work")
	SignMine := worksignCmd.Bool("mine", false, "Mine immediately on the same node")//新添加的mine值
	Signname :=  worksignCmd.String("name","","the name of auther")

	//作品版权授权需要的参数
	Buypassword := workshouquanCmd.String("buypassword","","the password of Buyer")
	Selladdress := workshouquanCmd.String("Selladdress","","the address of Seller")
	ShouquanTitle := workshouquanCmd.String("Title","","the Title of work")
	//Remarks := workshouquanCmd.String("Remarks","","the Remarks of shouquan")
	ShouquanMine := workshouquanCmd.Bool("mine", false, "Mine immediately on the same node")//新添加的mine值
	Buyname := workshouquanCmd.String("Buyname","","the Buyname of Buyer")

	//作品版权转权需要的参数
	ZBuypassword := workzhuanquanCmd.String("buypassword","","the password of Buyer")
	ZSelladdress := workzhuanquanCmd.String("Selladdress","","the address of Seller")
	ZhuanquanTitle := workzhuanquanCmd.String("Title","","the Title of work")
	//ZhuanquanRemarks := workzhuanquanCmd.String("Remarks","","the Remarks of Zhuanquan")
	NewAuthorizationmoney := workzhuanquanCmd.Int("NewAuthorizationmoney",0,"the NewAuthorizationmoney of Zhuanquan")
	NewTransactionmoney := workzhuanquanCmd.Int("NewTransactionmoney",0,"the NewTransactionmoney of Zhuanquan")
	ZhuanquanMine := workzhuanquanCmd.Bool("mine", false, "Mine immediately on the same node")//新添加的mine值
	ZBuyname := workzhuanquanCmd.String("Buyname","","the Buyname of Buyer")


	//通过密码找到用户地址
	password := getaddressCmd.String("password","","the password of wallet")
	getaddressname := getaddressCmd.String("getaddressname","","the name of wallet")

	//通过密码跟作品名找到区块链上对应的交易
	transactionpassword := gettransactionCmd.String("password","","the password of wallet")
	transactionname := gettransactionCmd.String("tranname","","the name of wallet")
	title := gettransactionCmd.String("title","","the title of transaction")
	transactiontype := gettransactionCmd.String("type","","the type of transaction")

	//展示作品的所有交易信息
	worktitle := getworkrecordCmd.String("title","","查询作品的名称")

	//查询用户信息需要的参数
	username := getuserCmd.String("name","","用户姓名")
	userpassword := getuserCmd.String("password","","用户密码")

	//创建作品文件
	workfile := createworkfileCmd.String("filename","","作品文件名")

	//作品版权溯源参数
	tracetitle := tracebackzhuanquanCmd.String("title","","作品名称")
	tracename := tracebackzhuanquanCmd.String("name","","用户姓名")
	tracepassword := tracebackzhuanquanCmd.String("password","","用户密码")

	//使用权溯源参数
	shouquantracetitle := tracebackshouquanCmd.String("title","","作品名称")
	shouquantracename := tracebackshouquanCmd.String("name","","用户姓名")
	shouquantracepassword := tracebackshouquanCmd.String("password","","用户密码")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printlastchain":
		err := printlastChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "作品版权签名":
		err := worksignCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "版权注册":
		err := banquanzhuceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "作品修改":
		err := rebanquanzhuceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "版权使用权购买":
		err := workshouquanCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "作品版权转让":
		err := workzhuanquanCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "启动节点":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getaddress":
		err := getaddressCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "gettransaction":
		err := gettransactionCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getallwork":
		err := getallworkCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getworkrecord":
		err := getworkrecordCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getuser":
		err := getuserCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createworkfile":
		err := createworkfileCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "tracebackzhuanquan":
		err := tracebackzhuanquanCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "tracebackshouquan":
		err := tracebackshouquanCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	//检查是哪个子命令并调用相关参数
	if getBalanceCmd.Parsed(){//解析Cmd.parsed是否被调用
		if *getBalancename == "" || *getBalancepassword == ""{//给了命令，没有参数仍然退出
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		address := cli.Getaddressbypassword(*getBalancepassword,*getBalancename,nodeID)
		cli.getBalance(address,nodeID)//查询某个地址的钱数
	}

	if createBlockchainCmd.Parsed(){
		if *createBlockchainAddress == "" {//给了命令，没有参数仍然退出
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress,nodeID)//创建一个区块
	}

	if createWalletCmd.Parsed() {
		if *WalletName == "" || *WalletRole == "" ||*WalletIdcard == "" || *WalletPassword == ""{
			createWalletCmd.Usage()
			os.Exit(1)
		}
		//cli.createWallet(nodeID)//创建一个钱包
		cli.CreateWalletuser(*WalletName,*WalletRole,*WalletIdcard,*WalletPassword,nodeID)//连带用户信息创建一个钱包
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)//打印这个节点的钱包的列表
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)//打印输出改节点拥有的整个区块链
	}

	if printlastChainCmd.Parsed() {
		cli.PrintLastChain(nodeID)//打印输出改节点拥有的整个区块链
	}

	if getaddressCmd.Parsed() {
		if *password == "" || *getaddressname == ""{
			getaddressCmd.Usage()
			os.Exit(1)
		}
		address := cli.Getaddressbypassword(*password,*getaddressname,nodeID)//打印输出改节点拥有的整个区块链
		fmt.Println(address)
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}

	if sendCmd.Parsed() {
		if *sendTo == "" || *sendAmount <= 0 ||*sendpassword ==""|| *sendname == ""{
			sendCmd.Usage()
			os.Exit(1)
		}
		sendFrom := cli.Getaddressbypassword(*sendpassword,*sendname,nodeID)
		cli.send(sendFrom, *sendTo, *sendAmount,*sendname ,nodeID, *sendMine)
	}

	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMiner)
	}

	if banquanzhuceCmd.Parsed() {//作品版权注册
		if *Userpassword == "" ||*Zhucename==""||*Title == "" || *Abstract == "" ||*ZhuceWorkfile ==""  || *Authorizationmoney <= 0|| *Transactionmoney <= 0{
			banquanzhuceCmd.Usage()
			os.Exit(1)
		}
		Useraddress := cli.Getaddressbypassword(*Userpassword,*Zhucename,nodeID)
		cli.Banquanzhuce(Useraddress,*Zhucename, *Title, *Abstract,*ZhuceWorkfile,*Authorizationmoney,*Transactionmoney,nodeID,*ZhuceMine)
	}

	if rebanquanzhuceCmd.Parsed() {//作品修改
		if *ReUserpassword == "" ||*ReZhucename==""||*ReTitle == "" || *ReAbstract == "" ||*ReZhuceWorkfile ==""  || *ReAuthorizationmoney <= 0|| *ReTransactionmoney <= 0{
			rebanquanzhuceCmd.Usage()
			os.Exit(1)
		}
		ReUseraddress := cli.Getaddressbypassword(*ReUserpassword,*ReZhucename,nodeID)
		cli.ReBanquanzhuce(ReUseraddress,*ReZhucename, *ReTitle, *ReAbstract,*ReZhuceWorkfile,*ReAuthorizationmoney,*ReTransactionmoney,nodeID,*ReZhuceMine)
	}

	if worksignCmd.Parsed() {//用户私钥对作品内容进行签名生成作品签名标签
		if *SignUserpassword == "" ||*SignTitle == "" || *Signworkfile == "" || *Signname=="" {
			worksignCmd.Usage()
			os.Exit(1)
		}
		SignUseraddress := cli.Getaddressbypassword(*SignUserpassword,*Signname,nodeID)
		//挖矿之前的时间
		//fmt.Println("输入交易之前的时间:")
		//t1 := time.Now().Unix()
		//时间戳转化为具体时间
		//fmt.Println(time.Unix(t1, 0).String())
		//n:= 200
		//for i:=1;i<=n;i++ {
		//	cli.Worksign(SignUseraddress, *SignTitle, *Signworkfile,*Signname,nodeID,*SignMine)
		//	//fmt.Println(i,"\n")
		//}
		cli.Worksign(SignUseraddress, *SignTitle, *Signworkfile,*Signname,nodeID,*SignMine)
		//挖矿之前的时间
		//fmt.Println("输入交易之后的时间:")
		//t2 := time.Now().Unix()
		//时间戳转化为具体时间
		//fmt.Println(time.Unix(t2, 0).String())
	}

	if workshouquanCmd.Parsed() {//作品使用权授权
		if *Buyname == "" ||*Buypassword == "" ||*ShouquanTitle == "" || *Selladdress == "" {
			workshouquanCmd.Usage()
			os.Exit(1)
		}
		Buyaddress := cli.Getaddressbypassword(*Buypassword,*Buyname,nodeID)
		cli.Shouquan(Buyaddress,*Buyname,*Selladdress, *ShouquanTitle,nodeID,*ShouquanMine)
	}
	if workzhuanquanCmd.Parsed() {//作品版权授权
		if *ZBuypassword == "" ||*ZBuyname == "" ||*ZhuanquanTitle == "" || *ZSelladdress == "" || *NewAuthorizationmoney <= 0|| *NewTransactionmoney <= 0{
			workzhuanquanCmd.Usage()
			os.Exit(1)
		}
		ZBuyaddress := cli.Getaddressbypassword(*ZBuypassword,*ZBuyname,nodeID)
		cli.Zhuanquan(ZBuyaddress,*ZBuyname ,*ZSelladdress, *ZhuanquanTitle ,*NewAuthorizationmoney,*NewTransactionmoney,nodeID,*ZhuanquanMine)
	}

	if gettransactionCmd.Parsed() {//找到对应的交易
		if *transactionpassword == "" ||*title == "" ||*transactiontype=="" || *transactionname == ""{
			gettransactionCmd.Usage()
			os.Exit(1)
		}
		//通过密码得到地址
		address := cli.Getaddressbypassword(*transactionpassword,*transactionname,nodeID)

		if *transactiontype =="授权记录" {
			//得到授权记录
			cli.Getshouquanbyaddressandtitle(address,*title,nodeID)
		} else if *transactiontype=="转权记录" {
			//得到转权记录
			cli.Getzhuanquanbyaddressandtitle(address,*title,nodeID)
		} else if *transactiontype=="认证记录"{
			//得到认证记录
		    cli.Getworkfromblock(address,*title,nodeID)
		}else if *transactiontype == "签名记录" {
			//得到签名记录
			cli.Getworkcontentbyaddressandtitlesign(address,*title,nodeID)

		}
	}

	if getallworkCmd.Parsed() {
		cli.Getallwork(nodeID)//展示所有的作品
	}

	if getworkrecordCmd.Parsed(){
		if *worktitle == "" {
			getworkrecordCmd.Usage()
			os.Exit(1)
		}
		cli.Getworkrecord(*worktitle,nodeID)//得到相关作品区块链上的信息
	}

	if getuserCmd.Parsed(){
		if *username == "" || *userpassword == ""{
			getuserCmd.Usage()
			os.Exit(1)
		}
		cli.Getuser(*username,*userpassword,nodeID)//得到相关作品区块链上的信息
	}

	if createworkfileCmd.Parsed(){
		if *workfile == "" {
			createworkfileCmd.Usage()
			os.Exit(1)
		}
		cli.CreateworkFile(*workfile,nodeID)//创建一个作品文件
	}

	if tracebackzhuanquanCmd.Parsed(){
		if *tracetitle == ""  || *tracename == "" || *tracepassword == ""{
			tracebackzhuanquanCmd.Usage()
			os.Exit(1)
		}

		address := cli.Getaddressbypassword(*tracepassword,*tracename,nodeID)
		cli.Traceback(*tracetitle,address,nodeID)//针对用户作品版权进行溯源查询
	}

	if tracebackshouquanCmd.Parsed(){
		if *shouquantracetitle == ""  || *shouquantracename == "" || *shouquantracepassword == ""{
			tracebackshouquanCmd.Usage()
			os.Exit(1)
		}

		address := cli.Getaddressbypassword(*shouquantracepassword,*shouquantracename,nodeID)
		cli.Tracebackshouquan(*shouquantracetitle,address,nodeID)//针对用户作品使用权溯源进行溯源查询
	}
}
