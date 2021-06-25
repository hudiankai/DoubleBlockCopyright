package core

import (
	"fmt"
	"strconv"
	"time"
)

func (cli *CLI) Getworkrecord(title,nodeID string) {
	//实例化一条链
	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()

	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体

	for {

		block := bci.Next()	//从顶端区块向前面的区块迭代

		for _,tx := range block.Transactions {

			if tx.WorkData.Title ==title {
				if tx.WorkData.Abstract == ""{
					fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
					fmt.Println("时间戳:",block.Timestamp)
					fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
					fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
					fmt.Println("区块高度:",block.Height)
					fmt.Println("随机数nonce:",block.Nonce)
					pow := NewProofOfWork(block)
					fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

					fmt.Println("交易附带的作品签名信息:")
					fmt.Println("作者ID:",tx.WorkData.Useridcard,",作者姓名:",tx.WorkData.Username,",作者地址:",tx.WorkData.Useraddress,
						",作品题目:",tx.WorkData.Title,",作品内容加密标签R:",tx.WorkData.Contentkeyrtext,
						"作品内容加密标签S",tx.WorkData.Contentkeystext,"时间戳:",tx.WorkData.Time,"\n")

				}else {
					fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
					fmt.Println("时间戳:",block.Timestamp)
					fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
					fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
					fmt.Println("区块高度:",block.Height)
					fmt.Println("随机数nonce:",block.Nonce)
					pow := NewProofOfWork(block)
					fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

					fmt.Println("交易附带的作品版权注册信息:")
					fmt.Println("作者ID:",tx.WorkData.Useridcard,",作者姓名:",tx.WorkData.Username,",作者地址:",tx.WorkData.Useraddress,
						",作品题目:",tx.WorkData.Title,",作品摘要:",tx.WorkData.Abstract,",作品内容加密标签R:",
						tx.WorkData.Contentkeyrtext,"作品内容加密标签S",tx.WorkData.Contentkeystext, ",作品版权授权费用:",tx.WorkData.Authorizationmoney,",作品版权转权费用:",
						tx.WorkData.Transactionmoney,"时间戳:",tx.WorkData.Time,"\n")
				}

			}

			if tx.ShouquanData.Title == title {
				fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
				fmt.Println("时间戳:",block.Timestamp)
				fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
				fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
				fmt.Println("区块高度:",block.Height)
				fmt.Println("随机数nonce:",block.Nonce)
				pow := NewProofOfWork(block)
				fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

				fmt.Println("交易附带的作品版权授权信息:")
				fmt.Println("买方地址:",tx.ShouquanData.Buyaddress,"卖方地址:",tx.ShouquanData.Selladdress,"授权作品名称:",tx.ShouquanData.Title,
					"交易transaction:",tx.ShouquanData.Transactionhash,"时间戳:",tx.ShouquanData.Time,"授权交易资金:",tx.ShouquanData.Money,"备注信息:",tx.ShouquanData.Remarks,"\n")
			}

			if tx.Zhuanquandata.WorkTital == title {
				fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
				fmt.Println("时间戳:",block.Timestamp)
				fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
				fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
				fmt.Println("区块高度:",block.Height)
				fmt.Println("随机数nonce:",block.Nonce)
				pow := NewProofOfWork(block)
				fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

				fmt.Println("交易附带的作品版权转权信息:")
				fmt.Println("买方地址:",tx.Zhuanquandata.Buyaddress,"卖方地址:",tx.Zhuanquandata.Selladdress,"转权作品名称:",tx.Zhuanquandata.WorkTital,
					"交易transaction:",tx.Zhuanquandata.Transactionhash,"时间戳:",tx.Zhuanquandata.Time,"转权交易资金:",tx.Zhuanquandata.ZhuanquanMoney,
					"新的授权价格:",tx.Zhuanquandata.Newauthorizationmoney,"新的转权价格:",tx.Zhuanquandata.Newtransactionmoney,"备注信息:",tx.Zhuanquandata.Remarks,"\n")

			}


		}
		fmt.Printf("\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}