package core

import (
	"fmt"
	"strconv"
	"time"
)

//列出所有的作品
func (cli *CLI) Getallwork(nodeID string) {

	var worknames []string

	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()

	b:= false
	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	fmt.Println("作品信息展示如下: \n")
	for {

		block := bci.Next()	//从顶端区块向前面的区块迭代

		for _,tx := range block.Transactions {
			if tx.WorkData != Worknil && tx.WorkData.Authorizationmoney!=0{
				for _,v :=range worknames {
					if v == tx.WorkData.Title {
						b =true
						break
					}
				}
				if b == false {

					fmt.Println("作品题目:",tx.WorkData.Title)
					fmt.Println("作者地址:",tx.WorkData.Useraddress)
					fmt.Println("作品摘要:",tx.WorkData.Abstract)
					fmt.Println("作品注册时间:",tx.WorkData.Time)
					fmt.Println("作品授权费用:",tx.WorkData.Authorizationmoney)
					fmt.Println("作品转权费用:",tx.WorkData.Transactionmoney)


					fmt.Println("对应的区块信息如下:")
					fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
					fmt.Println("时间戳:",block.Timestamp)
					fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
					fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
					fmt.Println("区块高度:",block.Height)
					pow := NewProofOfWork(block)
					boolen := pow.Validate()
					fmt.Printf("POW is %s\n",strconv.FormatBool(boolen))
					fmt.Println("\n")
					worknames = append(worknames,tx.WorkData.Title)
				}
				b = false
			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
//
//type Work struct {
//	Useridcard string//用户身份证号
//	Useraddress string//对应智能合约中的用户地址
//	Username string
//	Title string
//	Abstract string//作品内容摘要
//	WorkID int
//	//椭圆加密签名就是两个
//	Contentkeyrtext string//作品内容加密标签r,
//	Contentkeystext string//作品内容加密标签s,
//	Authorizationmoney int //版权授权费用
//	Transactionmoney int //版权转权费用
//	Transactionhash string
//	Time time.Time//时间戳
//}//与智能合约中的作品信息一一对应