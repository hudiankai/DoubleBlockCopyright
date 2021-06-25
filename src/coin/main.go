package main

import "blockchaincopyright/src/core"

func main() {
	//bc := core.NewBlockchain()
	//defer bc.Db.Close()//main函数结束之后关闭文件数据库


	cli := core.CLI{}//cli客户端的缩写
	cli.Run()

	//lixiaotongadd := core.GetaddressbyIDcard("19941004")
	//hudiankaiadd := core.GetaddressbyIDcard("19931213")
	//fmt.Println("hudiankai",hudiankaiadd,"lixiaotong",lixiaotongadd)

	//cli.PrintChain()
	//address := core.GetaddressbyIDcard("111111111111111")//通过身份证号查找钱包地址
	//fmt.Println(address)

	//创作带有用户信息的钱包测试
	//cli.CreateWalletuser("胡殿凯","用户一","19931213")//1LdbQNG8jBPkWb4KDcj5KjeEmNEF5EiAyh
	//cli.CreateWalletuser("李晓彤","用户一","19941004")//17eDT8ESB78vurBFgib7nGssBcrjrYJFXh


	//作品内容私钥签名测试
	//cli.Worksign(core.GetaddressbyIDcard("19931213"),"胡殿凯签名作品","用于测试")

	//作品版权注册测试,成功注册作品奖励10元
	//cli.GetBalance(core.GetaddressbyIDcard("19931213"))
	//cli.Banquanzhuce(core.GetaddressbyIDcard("19931213"),"胡殿凯的第二个作品","用于测试通过transactionHash找到交易内容","测试的用于测试通过transactionHash找到交易内容内容",2,2,5)
	//cli.GetBalance(core.GetaddressbyIDcard("19931213"))
	//cli.Banquanzhuce(lixiaotongadd,"李晓彤授权胡殿凯","用于测试通过transactionHash找到授权交易内容","测试的呵呵哈哈哈或或或或或或或或或或或或或或或内容",2,2,5)


	//作品版权授权测试,
	//cli.GetBalance(lixiaotongadd)
	//cli.GetBalance(hudiankaiadd)
	//cli.Shouquan(hudiankaiadd,lixiaotongadd,"李晓彤授权胡殿凯", "李晓彤请求胡殿凯授权他的作品<胡殿凯签名作品>,用于测试授权标签")
	//cli.GetBalance(lixiaotongadd)
	//cli.GetBalance(hudiankaiadd)

	//作品版权转权测试
	//cli.GetBalance("17eDT8ESB78vurBFgib7nGssBcrjrYJFXh")
	//cli.GetBalance("1LdbQNG8jBPkWb4KDcj5KjeEmNEF5EiAyh")
	//cli.Zhuanquan("1LdbQNG8jBPkWb4KDcj5KjeEmNEF5EiAyh","17eDT8ESB78vurBFgib7nGssBcrjrYJFXh","李晓彤的作品2",
	//	"胡殿凯购买李晓的作品的版权<李晓彤的作品2>用于测试查询交易hash",10,30)
	//cli.GetBalance("17eDT8ESB78vurBFgib7nGssBcrjrYJFXh")
	//cli.GetBalance("1LdbQNG8jBPkWb4KDcj5KjeEmNEF5EiAyh")


	//通过用户地址以及作品题目找到作品信息
	//cli.Getwalletbyaddress("1LdbQNG8jBPkWb4KDcj5KjeEmNEF5EiAyh")
	//cli.Getwalletbyaddress("17eDT8ESB78vurBFgib7nGssBcrjrYJFXh")
	//cli.Getworkbyaddressandtitle(core.GetaddressbyIDcard("19931213"),"胡殿凯的第二个作品")
	//cli.Getworkcontentbyaddressandtitle("1LdbQNG8jBPkWb4KDcj5KjeEmNEF5EiAyh","胡殿凯的第二个作品")
	//cli.Getworkbyaddressandtitle(core.GetaddressbyIDcard("19941004"),"李晓彤的作品2")
	//cli.Getworkcontentbyaddressandtitle(core.GetaddressbyIDcard("19941004"),"李晓彤的作品2")
	//cli.Getworkbyaddressandtitle("17eDT8ESB78vurBFgib7nGssBcrjrYJFXh","李晓彤的作品")
	//cli.Getworkcontentbyaddressandtitle("17eDT8ESB78vurBFgib7nGssBcrjrYJFXh","李晓彤的作品")
	//cli.Getshouquanbyaddressandtitle(hudiankaiadd,"李晓彤授权胡殿凯")

	//cli.PrintLastChain()//打印最新的一个区块
	//cli.PrintChain()//打印整个区块
	//通过hash值查找区块测试
	//cli.PrintChainbyhash("0004ee02e0f8c37a6cb98ca2768212801b86c541e2f666bed02cff626722bb21")

	//通过hash找到transaction
	//hash := cli.Getworkbyaddressandtitle(core.GetaddressbyIDcard("19941004"),"李晓彤的作品2").Transactionhash
	//cli.PrintTransactionbyhash(hash)
	//hash1 := cli.Getworkbyaddressandtitle(core.GetaddressbyIDcard("19931213"),"李晓彤的作品2").Transactionhash
	//cli.PrintTransactionbyhash(hash1)
	//hashc := cli.Getworkcontentbyaddressandtitle(core.GetaddressbyIDcard("19931213"),"胡殿凯签名作品").Transactionhash
	//cli.PrintTransactionbyhash(hashc)
	//cli.PrintTransactionbyhash(	cli.Getshouquanbyaddressandtitle(hudiankaiadd,"李晓彤授权胡殿凯").Transactionhash)
	//
	//shouquans := core.Shouquans{}
	//fmt.Println(shouquans.LoadshouquanFromFile())
	//workcontents := core.Workcontents{}
	//fmt.Println(workcontents.LoadworkcontentFromFile())

	//var docs = [][]byte{
	//	[]byte("中国人"),
	//	[]byte("中国"),
	//	[]byte("中国的人"),
	//}
	//
	//hashes := make([]uint64, len(docs))
	//for i, d := range docs {
	//	hashes[i] = simhash.Simhash(simhash.NewWordFeatureSet(d))
	//	fmt.Println( d,string(d), hashes[i])
	//}
	//
	//fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[1], simhash.Compare(hashes[0], hashes[1]))
	//fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[2], simhash.Compare(hashes[0], hashes[2]))
	//
	//st:=strconv.FormatUint(hashes[0],10)
	//in, _ := strconv.ParseUint(st, 10, 64)
	//fmt.Println(in)



}

