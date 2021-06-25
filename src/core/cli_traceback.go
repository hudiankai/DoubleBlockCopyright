package core

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func (cli *CLI) Traceback(title,address,nodeID string) {
	//实例化一条链
	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()
	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}

	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	traceaddress := address//用户地址(用户追溯用户解密)
	var bhash string//用于记录作品版权前一记录区块的hash

	for {
		if !ValidateAddress(traceaddress) {
			log.Panic("ERROR: BuyAddress is not valid")
		}
		block := bci.Next()	//从顶端区块向前面的区块迭代

		if traceaddress != address && string(block.Hash) != bhash{
			continue
		}//直接查询目标区块
		for _,tx := range block.Transactions {
			if tx.WorkData != Worknil && tx.WorkData.Authorizationmoney!=0{
				if  tx.WorkData.Title == title && traceaddress == tx.WorkData.Useraddress{
					if tx.Zhuanquandata.Buyaddress == traceaddress{//不是作品的原创作者,中间的版权购买者
						content := tx.WorkData.Content
						rtext := []byte(tx.WorkData.Contentkeyrtext)
						stext := []byte(tx.WorkData.Contentkeystext)
						userwallet := wallets.GetWallet(traceaddress)
						if VerifySignECCWork([]byte(content),rtext,stext,userwallet.PrivateKey.PublicKey){//作品签名验证失败
							log.Println("公钥解密成功")
							pow := NewProofOfWork(block)
							if pow.Validate(){//验证区块信息是否正确
								fmt.Println("区块信息验证成功.")
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
									"交易transaction:",tx.Hash(),"时间戳:",tx.Zhuanquandata.Time,"转权交易资金:",tx.Zhuanquandata.ZhuanquanMoney,
									"新的授权价格:",tx.Zhuanquandata.Newauthorizationmoney,"新的转权价格:",tx.Zhuanquandata.Newtransactionmoney,"\n")
								traceaddress = tx.Zhuanquandata.Selladdress//追溯上一个版权所有者的信息,验证版权
								bhash = tx.Zhuanquandata.Remarks
							}
						}
					} else if tx.Zhuanquandata == Zhuanquannil {//溯源到了作品的创作者
						content := tx.WorkData.Content
						rtext := []byte(tx.WorkData.Contentkeyrtext)
						stext := []byte(tx.WorkData.Contentkeystext)
						userwallet := wallets.GetWallet(traceaddress)
						if VerifySignECCWork([]byte(content),rtext,stext,userwallet.PrivateKey.PublicKey){//作品签名验证失败
							log.Println("公钥解密成功")
							pow := NewProofOfWork(block)
							if pow.Validate(){//验证区块信息是否正确
								fmt.Println("区块信息验证成功.")
								fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
								fmt.Println("时间戳:",block.Timestamp)
								fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
								fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
								fmt.Println("区块高度:",block.Height)
								fmt.Println("随机数nonce:",block.Nonce)
								pow := NewProofOfWork(block)
								fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

								fmt.Println("原创作者的版权注册信息:")
								fmt.Println("原创者地址:",tx.WorkData.Useraddress,"作品名称:",tx.WorkData.Title,"作品摘要:",tx.WorkData.Abstract,
									"交易transaction:",tx.Hash(),"时间戳:",tx.WorkData.Time,"备注信息:",tx.Zhuanquandata.Remarks,"\n")
								fmt.Println("本作品溯源完毕,本次作者是作品的原创作者")
								traceaddress = ""
							}
						}
					}
				}
			}

		}
		fmt.Printf("\n")

		if traceaddress == "" || len(block.PrevBlockHash) == 0  {
			break
		}
	}
}



func (cli *CLI) Tracebackshouquan(title,address,nodeID string) {//版权授权记录溯源
	//实例化一条链
	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()
	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}

	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	traceaddress := address//用户地址(用户追溯用户解密)
	var bhash string//用于记录作品版权前一记录区块的hash

	for {
		if !ValidateAddress(traceaddress) {
			log.Panic("ERROR: BuyAddress is not valid")
		}
		block := bci.Next()	//从顶端区块向前面的区块迭代

		if traceaddress == address && bhash == ""{
			for _,tx := range block.Transactions {
				if tx.ShouquanData.Title == title && traceaddress == tx.ShouquanData.Buyaddress {
					pow := NewProofOfWork(block)
					if pow.Validate() {
						fmt.Println("区块信息验证成功.")
						fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
						fmt.Println("时间戳:", block.Timestamp)
						fmt.Println("区块创建的时间(时间戳转化而来):", time.Unix(block.Timestamp, 0))
						fmt.Printf("PrevHash:%x\n", block.PrevBlockHash)
						fmt.Println("区块高度:", block.Height)
						fmt.Println("随机数nonce:", block.Nonce)
						fmt.Printf("POW is %s\n", strconv.FormatBool(pow.Validate()))

						fmt.Println("作品授权记录信息如下:")
						fmt.Println("作品题目:", tx.ShouquanData.Title, "授权方地址", tx.ShouquanData.Selladdress, "\n",
							"授权时间:", tx.ShouquanData.Time, "授权价格:", tx.ShouquanData.Money, "\n")
						traceaddress = tx.ShouquanData.Selladdress//追溯上一个版权所有者的信息,验证版权
						bhash = tx.ShouquanData.Remarks//授权方记录所在的区块hash
					}
				}
			}
		}else if string(block.Hash) == bhash {
			for _,tx := range block.Transactions {
				if tx.WorkData != Worknil && tx.WorkData.Authorizationmoney!=0{
					if  tx.WorkData.Title == title && traceaddress == tx.WorkData.Useraddress{
						if tx.Zhuanquandata.Buyaddress == traceaddress{//不是作品的原创作者,中间的版权购买者
							content := tx.WorkData.Content
							rtext := []byte(tx.WorkData.Contentkeyrtext)
							stext := []byte(tx.WorkData.Contentkeystext)
							userwallet := wallets.GetWallet(traceaddress)
							if VerifySignECCWork([]byte(content),rtext,stext,userwallet.PrivateKey.PublicKey){//作品签名验证失败
								log.Println("公钥解密成功")
								pow := NewProofOfWork(block)
								if pow.Validate(){//验证区块信息是否正确
									fmt.Println("区块信息验证成功.")
									fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
									fmt.Println("时间戳:",block.Timestamp)
									fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
									fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
									fmt.Println("区块高度:",block.Height)
									fmt.Println("随机数nonce:",block.Nonce)
									fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

									fmt.Println("作品版权转权信息:")
									fmt.Println("买方地址:",tx.Zhuanquandata.Buyaddress,"卖方地址:",tx.Zhuanquandata.Selladdress,"转权作品名称:",tx.Zhuanquandata.WorkTital,
										"交易transaction:",tx.Hash(),"时间戳:",tx.Zhuanquandata.Time,"转权交易资金:",tx.Zhuanquandata.ZhuanquanMoney,
										"新的授权价格:",tx.Zhuanquandata.Newauthorizationmoney,"新的转权价格:",tx.Zhuanquandata.Newtransactionmoney,"\n")
									traceaddress = tx.Zhuanquandata.Selladdress//追溯上一个版权所有者的信息,验证版权
									bhash = tx.Zhuanquandata.Remarks
								}
							}
						} else if tx.Zhuanquandata == Zhuanquannil {//溯源到了作品的创作者
							content := tx.WorkData.Content
							rtext := []byte(tx.WorkData.Contentkeyrtext)
							stext := []byte(tx.WorkData.Contentkeystext)
							userwallet := wallets.GetWallet(traceaddress)
							if VerifySignECCWork([]byte(content),rtext,stext,userwallet.PrivateKey.PublicKey){//作品签名验证失败
								log.Println("公钥解密成功")
								pow := NewProofOfWork(block)
								if pow.Validate(){//验证区块信息是否正确
									fmt.Println("区块信息验证成功.")
									fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
									fmt.Println("时间戳:",block.Timestamp)
									fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
									fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
									fmt.Println("区块高度:",block.Height)
									fmt.Println("随机数nonce:",block.Nonce)
									fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

									fmt.Println("原创作者的版权注册信息:")
									fmt.Println("原创者地址:",tx.WorkData.Useraddress,"作品名称:",tx.WorkData.Title,"作品摘要:",tx.WorkData.Abstract,
										"交易transaction:",tx.Hash(),"时间戳:",tx.WorkData.Time,"备注信息:",tx.Zhuanquandata.Remarks,"\n")
									fmt.Println("本作品溯源完毕,本次作者是作品的原创作者")
									traceaddress = ""
								}
							}
						}
					}
				}

			}
		}

		//fmt.Printf("\n")

		if traceaddress == "" || len(block.PrevBlockHash) == 0  {
			break
		}
	}
}