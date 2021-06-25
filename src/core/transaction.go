package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var subsidy int = 50  //挖矿奖励
var WorkMoney int = 10 //作品版权注册奖励
/*创建一个交易的数据结构，交易是由交易ID、交易输入、交易输出组成的,
一个交易有多个输入和多个输出，所以这里的交易输入和输出应该是切片类型的
*/
type Transaction struct {
	ID		[]byte//整个结构体的数据序列化然后hash运算得到的hash
	Vin		[]TXInput
	Vout	[]TXOutput
	WorkData Work
	ShouquanData Shouquan
	Zhuanquandata Zhuanquan
}

var Worknil = Work{}
var Shouquannil = Shouquan{}
var Zhuanquannil = Zhuanquan{}


//交易输出
type TXOutput struct {
	Value			int	//输出的值（可以理解为金额）
	PubkeyHash 		[]byte //存储“哈希”后的公钥，这里的哈希不是单纯的sha256
}

//交易输入
type TXInput struct {
	Txid 		[]byte //引用的之前交易的ID
	Vout		int 	//该笔交易输出的索引,一笔交易可能有多个输出
	Signature 	[]byte //私钥签名脚本(私钥签名生成的r ,s数字的链接)
	PubKey 		[]byte // 公钥，这里的公钥是正真的公钥
}

//创建一个结构体，用于表示TXOutput集
type TXOutputs struct {
	Outputs []TXOutput
}
//序列化此集合


//现在我们来创建一个这样的coinbase挖矿输出,to代表此输出奖励给谁，一般都是矿工地址，data是交易附带的信息
func NewCoinbaseTX(to,data string) *Transaction {

	if data == "" {
		data = fmt.Sprintf("奖励给 '%s'",to)
	}
	//此交易中的交易输入,没有交易输入信息
	txin := TXInput{[]byte{},-1,nil,[]byte(data)}
	//交易输出,subsidy为奖励矿工的币的数量
	txout := NewTXOutput(subsidy,to)

	//组成交易
	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{*txout},Worknil,Shouquannil,Zhuanquannil}

	//设置该交易的ID
	tx.ID = tx.Hash()
	return &tx
}

//判断是否为coinbase交易
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}//

////设置交易ID，交易ID是序列化tx后再哈希
//返回一个序列化后的交易
func (tx Transaction) Serialize() []byte {
	var encoder bytes.Buffer

	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoder.Bytes()
}
//返回交易的哈希值,
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())//先将transaction序列化,然后再将其进行hash运算

	return hash[:]
}


//Useskey检查地址是否可用
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)//要求是哈希后的公钥

	return bytes.Compare(lockingHash,pubKeyHash) == 0
}


//锁定交易输出到固定的地址，代表该输出只能由指定的地址引用
func (out *TXOutput) Lock(address []byte) {//对输出进行签名(锁定)
	pubKeyHash := Base58Decode(address)//将地址进行解码
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]//提取出公钥哈希
	out.PubkeyHash = pubKeyHash
}

//判断输入的公钥"哈希"能否解锁该交易输出
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubkeyHash,pubKeyHash) == 0  //如果两个字节切片相同则返回0
}

//创建一个新的交易输出
func NewTXOutput(value int,address string) *TXOutput {
	txo := &TXOutput{value,nil}
	txo.Lock([]byte(address))

	return txo
}



//对交易签名,
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey,prevTXs map[string]Transaction) {
	//coinbase交易没有输入所以不需要进行签名
	if tx.IsCoinbase() {
		return
	}

	for _,vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()//创建一个副本

	//迭代副本中的每一个输入
	for inID,vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]//是一个交易
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubkeyHash
		//Hash方法对交易进行序列化.使用SHA-256算法进行哈希,最后结果就是我们要签名的数据
		txCopy.ID = txCopy.Hash()
		//获取完哈希,我们应该重置Pubkey字段,以便于它不会影响后面的迭代
		txCopy.Vin[inID].PubKey = nil

		//通过privkey对txcopy.ID进行签名
		r,s,err := ecdsa.Sign(rand.Reader,&privKey,txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		//一个 ECDSA签名就是一对数字,我们把这对数字链接起来,并存储在输入的signature字段
		signature := append(r.Bytes(),s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}

}

///验证交易中输入的签名
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _,vin := range tx.Vin { //遍历输入交易，如果发现输入交易引用的上一交易的ID不存在，则Panic
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy() //修剪后的副本

	curve := elliptic.P256() //椭圆曲线实例

	for inID,vin := range tx.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil //双重验证
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.Vout].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil
		//解包存储在TXINPUT.Signature和TXInput.pubkey中的值,一个签名就是一对数字,一个公钥就是一堆坐标


		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		//使用从输入提取的公钥创建了一个ecdsa.Publickey,通过传入输出中提取的签名执行了 ecdsa.Verify
		//如果所有的输入都被验证，返回true；如果有任何一个交易失败，返回 false
		rawPubKey := ecdsa.PublicKey{curve,&x,&y}
		if ecdsa.Verify(&rawPubKey,txCopy.ID,&r,&s) == false {
			return false
		}
	}
	return true
}


//创建在签名中修剪后的交易副本,之所以要这个副本是因为简化了输入交易本身的签名和公钥
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _,vin := range tx.Vin {
		inputs = append(inputs,TXInput{vin.Txid,vin.Vout,nil,nil})
	}

	for _,vout := range tx.Vout {
		outputs = append(outputs,TXOutput{vout.Value,vout.PubkeyHash})
	}

	txCopy := Transaction{tx.ID,inputs,outputs,tx.WorkData,tx.ShouquanData,tx.Zhuanquandata}

	return txCopy
}

//把交易转换成我们能正常读的形式
func (tx Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("--Transaction %x:", tx.ID))//将Transaction的id写入到lines字符串数组中
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf(" -InputID %d:", i))
		lines = append(lines, fmt.Sprintf("  TXID: %x", input.Txid))
		lines = append(lines, fmt.Sprintf("  Out:  %d", input.Vout))
		lines = append(lines, fmt.Sprintf("  Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("  PubKey:%x", input.PubKey))
	}
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf(" -OutputID %d:", i))
		lines = append(lines, fmt.Sprintf("  Value: %d", output.Value))
		lines = append(lines, fmt.Sprintf("  Script: %x", output.PubkeyHash))
	}
	return strings.Join(lines,"\n")

}


//反序列化,由字节切片到交易输出的集合
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}

// DeserializeTransaction反序列化一个事务
func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

//序列化此集合
func(outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

//创建一笔新的不同交易,确保有足够的余额才能作为输出
func NewUTXOTransaction(from,to,sendname ,nodeID string,amount int,UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}
	_wallet := wallets.GetWallet(from)//通过发送地址得到发送方的钱包
	if sendname != _wallet.Userinformation.Name {
		log.Panic("用户支付密码错误,请重新输入密码.")
	}
	pubKeyHash := HashPubKey(_wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)//查询并返回被用于这次花费的输出，
	// 找到的输出的总额要刚好大于要花费的输入额,之前的交易的to是这次要查询账户
	if acc < amount {
		log.Panic("ERROR:Not enough tokens...")
	}

	//通过validOutputs里面的数据来放入建立一个输入列表
	for txid,outs := range validOutputs {
		//反序列化得到txID
		txID,err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//遍历输出outs切片,得到TXInput里的Vout字段值
		for _,out := range outs {
			input := TXInput{txID,out,nil,_wallet.PublicKey}
			inputs = append(inputs,input)
		}
	}
	//建立一个输出列表
	outputs = append(outputs,*NewTXOutput(amount,to))

	//找零
	if acc > amount {
		outputs = append(outputs,*NewTXOutput(acc - amount,from)) //相当于找零
	}

	tx := Transaction{nil,inputs,outputs,Worknil,Shouquannil,Zhuanquannil}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, _wallet.PrivateKey)

	return &tx
}


//创建一个包含作品版权注册的transaction
func NewBanquanTransaction(Useraddress,name,Title,Abstract,Content,nodeID ,simhash string,Authorizationmoney,Transactionmoney int,UTXOSet *UTXOSet) *Transaction {
	//validOutputs是一个存放要用到的未花费输出的交易/输出的map
	//acc,validOutputs := bc.FindSpendableOutputs(from,amount)
	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}
	_wallet := wallets.GetWallet(Useraddress)//通过发送地址得到发送方的钱包
	if name != _wallet.Userinformation.Name {
		log.Panic("用户名与密码不匹配!")
	}

	data := fmt.Sprintf("注册版权,奖励给 '%s'",_wallet.Userinformation.Name)
	//此交易中的交易输入,没有交易输入信息
	txin := TXInput{[]byte{},-1,nil,[]byte(data)}

	txout := NewTXOutput(WorkMoney,Useraddress)//交易输出,每注册一个作品奖励10个币,WorkMoney为10

	rtext,stext:=SignECCWork([]byte(Content),_wallet.PrivateKey)//椭圆曲线密钥对  私钥加密作品内容,生成作品私钥标签
	//把字节切片的签名结果转换成字符串的形式,
	contentrkey := string(rtext)
	contentskey := string(stext)

	workdata := Work{_wallet.Userinformation.Idcard,Useraddress,_wallet.Userinformation.Name,Title,
		Abstract,Content,contentrkey,contentskey,Authorizationmoney,
		Transactionmoney,"",time.Now(),simhash}
	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{*txout},workdata,Shouquannil,Zhuanquannil}


	tx.ID = tx.Hash()
    transactionhash := string(tx.ID)
	//UTXOSet.Blockchain.SignTransactionWork(&tx, _wallet.PrivateKey)

	works, _ := NewWorks()//返回一个结构体,只包含一个string到works的映射,并且加载work文件中的作品信息,
	biaoqian := works.CreateWork(Useraddress, Title,Abstract,Content,contentrkey,contentskey,nodeID,Authorizationmoney,Transactionmoney,transactionhash,simhash)
	works.SaveworkToFile()//
	//fmt.Println("作品标签",biaoqian,"作品",works.Works[biaoqian])

	//创建workcontent,讲作品内容信息存储在workcontent.date文件中
	workcontents,_ := NewWorkcontents()
	newworkcontent := Workcontent{_wallet.Userinformation.Idcard,Useraddress,_wallet.Userinformation.Name,Title,
		Abstract,Content,"",contentrkey,contentskey,Authorizationmoney,
		Transactionmoney,transactionhash,time.Now()}
	workcontents.Workcontents[biaoqian] = &newworkcontent
	workcontents.SaveworkcontentToFile()

	return &tx
}

//创建一个包含私钥对作品内容签名的transaction
func NewWorksignTransaction(SignUseraddress,SignTitle,SignContent,name,nodeID string,UTXOSet *UTXOSet) *Transaction {
	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}
	_wallet := wallets.GetWallet(SignUseraddress)//通过地址得到发送方的钱包
	if name != _wallet.Userinformation.Name {
		log.Panic("用户名与密码不匹配!")
	}

	data := fmt.Sprintf("作品版权签名登记,奖励给 '%s'",_wallet.Userinformation.Name)
	//此交易中的交易输入,没有交易输入信息
	txin := TXInput{[]byte{},-1,nil,[]byte(data)}	//交易输出,subsidy为奖励矿工的币的数量
	txout := NewTXOutput(0,SignUseraddress)//作品内容私钥签名登记不给于奖励

	rtext,stext:=SignECCWork([]byte(SignContent),_wallet.PrivateKey)
	//把字节切片的签名结果转换成字符串的形式,
	contentrkey := string(rtext)
	contentskey := string(stext)
	workdata := Work{_wallet.Userinformation.Idcard,SignUseraddress,_wallet.Userinformation.Name,SignTitle,
		"作品"+SignTitle+"签名登记",SignContent,contentrkey,contentskey,0,0,"",time.Now(),""}
	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{*txout},workdata,Shouquannil,Zhuanquannil}
	//tx.SetID()


	tx.ID = tx.Hash()
	transactionhash := string(tx.ID)
	//UTXOSet.Blockchain.SignTransactionWork(&tx, _wallet.PrivateKey)
	workdata.Transactionhash = string(tx.ID)
	//创建workcontent,讲作品内容信息存储在workcontent.date文件中
	workcontents,_ := NewWorkcontents()
	newworkcontent := Workcontent{_wallet.Userinformation.Idcard,SignUseraddress,_wallet.Userinformation.Name,SignTitle,
		"未完成作品登记",SignContent,"",contentrkey,contentskey,0,0,transactionhash,time.Now()}
	signbiaoqian := fmt.Sprintf("%s",newworkcontent.GetSHA256Signbiaoqian([]byte(SignUseraddress+SignTitle+"sign")))//获得每个作品对应的标签
	workcontents.Workcontents[signbiaoqian] = &newworkcontent
	workcontents.SaveworkcontentToFile()

	return &tx
}

//创建一个包含作品版权授权的transaction
func NewShouquanTransaction(buyaddress,buyname,selladdress, ShouquanTitle,nodeID string,UTXOSet *UTXOSet,work Work,bhash []byte) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}
	shouquans,_ := NewShouquans()//创建一个授权记录,读取shouquan.data文件中的信息

	buywallet := wallets.GetWallet(buyaddress)//通过发送地址得到发送方的钱包
	if buyname !=buywallet.Userinformation.Name {
		log.Panic("输入密码与用户不匹配!")
	}
	sellwallet := wallets.GetWallet(selladdress)//通过发送地址得到购买方的钱包
	pubKeyHash := HashPubKey(buywallet.PublicKey)

	//区块链上的作品数据
	worktransaction := work
	if worktransaction == Worknil {
		log.Panic("区块链中没有对应的作品")
	}
	ShouquanMoney := worktransaction.Authorizationmoney

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, ShouquanMoney)//查询并返回被用于这次花费的输出，
	// 找到的输出的总额要刚好大于购买版权所需的费用,之前的交易的to是这次要查询账户
	if acc < ShouquanMoney {
		log.Panic("ERROR:没有足够的钱购买该作品的版权使用权.")
	}
	//通过validOutputs里面的数据来放入建立一个输入列表
	for txid,outs := range validOutputs {
		//反序列化得到txID
		txID,err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//遍历输出outs切片,得到TXInput里的Vout字段值
		for _,out := range outs {
			//input := transaction.TXInput{txID,out,from}
			input := TXInput{txID,out,nil,buywallet.PublicKey}
			inputs = append(inputs,input)
		}
	}
	//建立一个输出列表
	outputs = append(outputs,*NewTXOutput(ShouquanMoney,selladdress))
	if acc > ShouquanMoney {
		outputs = append(outputs,*NewTXOutput(acc - ShouquanMoney,buyaddress)) //相当于找零
	}

	shouquandate := Shouquan{buyaddress,selladdress,ShouquanTitle,ShouquanMoney,string(bhash),"",
		time.Now(),worktransaction.Content}



	//验证授权方是否是现在版权的所有者
	rtext := []byte(worktransaction.Contentkeyrtext)
	stext := []byte(worktransaction.Contentkeystext)

	if !VerifySignECCWork([]byte(worktransaction.Content),rtext,stext,sellwallet.PrivateKey.PublicKey){//作品签名验证失败
		log.Println("Error:授权方公钥解密失败.")
	}

	tx := Transaction{nil,inputs,outputs,Worknil,shouquandate,Zhuanquannil}
	//tx.SetID()
	tx.ID = tx.Hash()
	//计算授权记录对应的记录标签Hash(买方地址+作品题目),然后存储记录
	shouquandate.Transactionhash = string(tx.ID)
	shouquanbiaoqian := shouquandate.GetSHA256Shouquanbiaoqian([]byte(buyaddress+ShouquanTitle))
	shouquans.Shouquans[shouquanbiaoqian] = &shouquandate
	shouquans.SaveshouquanToFile()
	UTXOSet.Blockchain.SignTransaction(&tx, buywallet.PrivateKey)

	return &tx
}


//创建一个包含作品版权转权的transaction
func NewZhuanquanTransaction(ZBuyaddress,ZBuyname ,ZSelladdress, ZhuanquanTitle,nodeID string,NewAuthorizationmoney,NewTransactionmoney int, UTXOSet *UTXOSet,worktra Work,bhash []byte) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}
	buywallet := wallets.GetWallet(ZBuyaddress)//通过发送地址得到购买方的钱包
	if ZBuyname !=buywallet.Userinformation.Name {
		log.Panic("输入密码与用户不匹配!")
	}
	sellwallet := wallets.GetWallet(ZSelladdress)//通过发送地址得到卖方的钱包

	pubKeyHash := HashPubKey(buywallet.PublicKey)

	works, _ := NewWorks()//返回一个结构体,只包含一个string到works的映射,并且加载work文件中的作品信息,
	work := Work{}//定义一个空的work用于调用GetSha256biaoqian函数获得作品的标签
	biaoqian := work.GetSHA256biaoqian([]byte(ZSelladdress+ZhuanquanTitle))//通过用户地址以及作品标题读取到作品的标签,以便找到作品规定的转权费用

	//区块链上的作品数据
	worktransaction := worktra
	if worktransaction == Worknil {
		log.Panic("区块链中没有对应的作品")
	}
	ZhuanquanMoney := worktransaction.Transactionmoney

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, ZhuanquanMoney)//查询并返回被用于这次花费的输出，
	// 找到的输出的总额要刚好大于购买版权所需的费用,之前的交易的to是这次要查询账户
	if acc < ZhuanquanMoney {
		log.Panic("ERROR:没有足够的钱购买该作品的版权所有权.")
	}
	//通过validOutputs里面的数据来放入建立一个输入列表
	for txid,outs := range validOutputs {
		//反序列化得到txID
		txID,err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//遍历输出outs切片,得到TXInput里的Vout字段值
		for _,out := range outs {
			//input := transaction.TXInput{txID,out,from}
			input := TXInput{txID,out,nil,buywallet.PublicKey}
			inputs = append(inputs,input)
		}
	}
	//建立一个输出列表
	outputs = append(outputs,*NewTXOutput(ZhuanquanMoney,ZSelladdress))
	if acc > ZhuanquanMoney {
		outputs = append(outputs,*NewTXOutput(acc - ZhuanquanMoney,ZBuyaddress)) //相当于找零
	}

	zhuanquandate := Zhuanquan{ZBuyaddress,ZSelladdress,ZhuanquanTitle,ZhuanquanMoney,
		time.Now(), NewAuthorizationmoney,NewTransactionmoney,"",string(bhash)}


	workcontents,error := NewWorkcontents()//
	if err != nil {
		log.Panic(error)
	}
	content := worktransaction.Content
	//验证卖方是否是现在版权的所有者
	oldrtext := []byte(worktransaction.Contentkeyrtext)
	oldstext := []byte(worktransaction.Contentkeystext)

	if !VerifySignECCWork([]byte(content),oldrtext,oldstext,sellwallet.PrivateKey.PublicKey){//作品签名验证失败
		log.Println("Error:卖方公钥解密失败.")
	}
	//有关于作品转权之后作品归属加密的调整
	newrtext,newstext:=SignECCWork([]byte(content),buywallet.PrivateKey)//椭圆曲线密钥对私钥加密作品内容,生成作品私钥标签
	//把字节切片的签名结果转换成字符串的形式,
	newcontentrkey := string(newrtext)
	newcontentskey := string(newstext)
	//fmt.Println(contentrkey,contentskey)

	//为新的版权所有者创作一个对应的作品信息
	workdata := Work{buywallet.Userinformation.Idcard,ZBuyaddress,buywallet.Userinformation.Name,worktransaction.Title,
		worktransaction.Abstract,content,newcontentrkey,newcontentskey,
		NewAuthorizationmoney,NewTransactionmoney,"",time.Now(),""}

	tx := Transaction{nil,inputs,outputs,workdata,Shouquannil,zhuanquandate}
	//tx.SetID()
	tx.ID = tx.Hash()
	transactionhash := string(tx.ID)
	workdata.Transactionhash = string(tx.ID)
	zhuanquandate.Transactionhash = string(tx.ID)
	UTXOSet.Blockchain.SignTransaction(&tx, buywallet.PrivateKey)

	//通过hash加密得到新的版权所有者对应的购买作品的版权标签
	newbiaoqian := work.GetSHA256biaoqian([]byte(ZBuyaddress+ZhuanquanTitle))
	//创建workcontent,将新的版权所有者对应的购买的作品的内容信息存储在workcontent.date文件中
	//workcontents,_ := NewWorkcontents()
	newworkcontent := Workcontent{buywallet.Userinformation.Idcard,ZBuyaddress,buywallet.Userinformation.Name,worktransaction.Title,
		worktransaction.Abstract,content,"",newcontentrkey,newcontentskey,NewAuthorizationmoney,
		NewTransactionmoney,transactionhash,time.Now()}
	workcontents.Workcontents[newbiaoqian] = &newworkcontent
	//作品转权之后,原版权所有者的作品信息修改以及新版权所有者的作品信息的添加
	works.Works[newbiaoqian] = &workdata//新的版权所有者生成新的Work信息存储在work.data文件中
	works.Works[newbiaoqian].Transactionhash = transactionhash//给新生成的买方的work.data复制交易hash
	works.Works[biaoqian].Contentkeyrtext =""//原作者不再拥有版权,故他的r签名秘钥修改为"",不能再次验证
	works.Works[biaoqian].Contentkeystext = ""//原作者不再拥有版权,故他的s签名秘钥修改为"",不能再次验证
	workcontents.Workcontents[biaoqian].Contentkeyrtext =""//原作者不再拥有版权,故他的r签名秘钥修改为"",不能再次验证
	workcontents.Workcontents[biaoqian].Contentkeystext = ""//原作者不再拥有版权,故他的s签名秘钥修改为"",不能再次验证
	workcontents.SaveworkcontentToFile()
	works.SaveworkToFile()//修改作品文件中的信息之后,将新的信息保存到work.data文件中

	zhuanquans,error := NewZhuanquans()//转权记录
	if err != nil {
		log.Panic(error)
	}
	zhuanquanbiaoqian := zhuanquandate.GetSHA256Zhuanquanbiaoqian([]byte(ZBuyaddress+ZhuanquanTitle))
	zhuanquans.Zhuanquans[zhuanquanbiaoqian] = &zhuanquandate
	zhuanquans.SavezhuanquanToFile()//转权记录写入数据库

	//fmt.Println("作品新的标签",newbiaoqian,"作品",works.Works[newbiaoqian])

	return &tx
}//work.date文件中的信息还没有更改



func (cli *CLI)Getworkfromblock(useraddress,title,nodeID string)  {
	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()

	b:= false
	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	for {

		block := bci.Next()	//从顶端区块向前面的区块迭代

		for _,tx := range block.Transactions {
			if tx.WorkData.Title == title && useraddress == tx.WorkData.Useraddress && tx.WorkData.Authorizationmoney != 0{

				//数字签名解密验证
				rtext := []byte(tx.WorkData.Contentkeyrtext)
				stext := []byte(tx.WorkData.Contentkeystext)
				wallet := Getwalletbyaddress(useraddress,nodeID)

				if !VerifySignECCWork([]byte(tx.WorkData.Content),rtext,stext,wallet.PrivateKey.PublicKey){//作品签名验证失败
					log.Println("Error:数字签名解密成功!.")
				}


				fmt.Println("认证作品信息如下:")
				fmt.Println("作品题目:",tx.WorkData.Title, "认证时间:",tx.WorkData.Time,"\n",
					"作品描述:",tx.WorkData.Abstract,"\n","作品内容:",tx.WorkData.Content,"\n")

				fmt.Println("对应的区块信息如下:")
				fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
				fmt.Println("时间戳:",block.Timestamp)
				fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
				fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
				fmt.Println("区块高度:",block.Height)
				pow := NewProofOfWork(block)
				boolen := pow.Validate()
				fmt.Printf("POW is %s\n",strconv.FormatBool(boolen))
				b = true
			}
			if b {
				break
			}
		}
		if b {
			break
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

//通过用户地址以及作品题目找到作品内容信息(认证之前的作品记录)
func (cli *CLI)Getworkcontentbyaddressandtitlesign(useraddress string,title,nodeID string)   {
	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()

	b := false
	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	for {


		block := bci.Next()	//从顶端区块向前面的区块迭代

		for _,tx := range block.Transactions {
			if tx.WorkData.Title == title && useraddress == tx.WorkData.Useraddress && tx.WorkData.Authorizationmoney == 0 && tx.WorkData.Transactionmoney == 0{
				fmt.Println("作品签名信息如下:")
				fmt.Println("作品题目:",tx.WorkData.Title, "签名时间:",tx.WorkData.Time,"\n",
					"作品描述:",tx.WorkData.Abstract,"\n","作品内容:",tx.WorkData.Content,"\n")

				fmt.Println("对应的区块信息如下:")
				fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
				fmt.Println("时间戳:",block.Timestamp)
				fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
				fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
				fmt.Println("区块高度:",block.Height)
				pow := NewProofOfWork(block)
				boolen := pow.Validate()
				fmt.Printf("POW is %s\n",strconv.FormatBool(boolen))
				b =true
			}

			if b {
				break
			}
		}

		if b {
			break
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}


//通过用户地址以及作品题目找到作品信息
func (cli *CLI)Getzhuanquanbyaddressandtitle(useraddress string,title ,nodeID string) {
	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()

	b := false
	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	for {

		block := bci.Next()	//从顶端区块向前面的区块迭代

		for _,tx := range block.Transactions {
			if tx.Zhuanquandata.WorkTital == title && useraddress == tx.Zhuanquandata.Buyaddress{
				fmt.Println("作品转权信息如下:")
				fmt.Println("作品题目:",tx.Zhuanquandata.WorkTital,"转权方地址",tx.Zhuanquandata.Selladdress,"\n",
					"转权时间:",tx.Zhuanquandata.Time,"转权价格:",tx.Zhuanquandata.ZhuanquanMoney,"\n","作品内容:",tx.WorkData.Content,"\n")


				fmt.Println("对应的区块信息如下:")
				fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
				fmt.Println("时间戳:",block.Timestamp)
				fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
				fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
				fmt.Println("区块高度:",block.Height)
				pow := NewProofOfWork(block)
				boolen := pow.Validate()
				fmt.Printf("POW is %s\n",strconv.FormatBool(boolen))
				b =true
			}
			if b {
				break
			}
		}

		if b {
			break
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

//通过用户地址以及作品题目找到作品信息
func (cli *CLI)Getshouquanbyaddressandtitle(useraddress string,title ,nodeID string) {
	bc := NewBlockchain(nodeID)  //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()

	b := false
	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()//一个存储最近区块hash和DB的结构体
	for {

		block := bci.Next()	//从顶端区块向前面的区块迭代

		for _,tx := range block.Transactions {
			if tx.ShouquanData.Title == title && useraddress == tx.ShouquanData.Buyaddress{
				fmt.Println("授权作品相关信息如下:")
				fmt.Println("作品题目:",tx.ShouquanData.Title,"授权方地址",tx.ShouquanData.Selladdress,"\n",
					"授权时间:",tx.ShouquanData.Time,"授权价格:",tx.ShouquanData.Money,"\n","作品内容:",tx.ShouquanData.Workcontent,"\n")

				fmt.Println("对应的区块信息如下:")
				fmt.Printf("------======= 区块Hash %x ============\n", block.Hash)
				fmt.Println("时间戳:",block.Timestamp)
				fmt.Println("区块创建的时间(时间戳转化而来):",time.Unix(block.Timestamp,0))
				fmt.Printf("PrevHash:%x\n",block.PrevBlockHash)
				fmt.Println("区块高度:",block.Height)
				pow := NewProofOfWork(block)
				boolen := pow.Validate()
				fmt.Printf("POW is %s\n",strconv.FormatBool(boolen))
				b = true
			}
			if b {
				break
			}
		}

		if b {
			break
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func  (bk Block)numbertx() int {
	i:=0
	for tx := range bk.Transactions {
		i++
		if i== 13333{
			fmt.Println(tx)
		}
	}
	return i
}