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
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

const workFile = "work.date"//存储写入区块链中的作品的版权等级记录以及作品转权记录
const workcontentFile = "workcontent.date"//存储写入区块链中的作品的具体内容等消息


type Work struct {
	Useridcard string//用户身份证号
	Useraddress string//对应智能合约中的用户地址
	Username string
	Title string
	Abstract string//作品内容摘要
	Content string//作品具体内容
	//椭圆加密签名就是两个
	Contentkeyrtext string//作品内容加密标签r,
	Contentkeystext string//作品内容加密标签s,
	Authorizationmoney int //版权授权费用
	Transactionmoney int //版权转权费用
	Transactionhash string
	Time time.Time//时间戳
	Contentsimhash string
}//与智能合约中的作品信息一一对应

type Workcontent struct {
	Useridcard string//用户身份证号
	Useraddress string//对应智能合约中的用户地址
	Username string
	Title string
	Abstract string//作品内容摘要
	Content string//作品具体内容
	Remarks string //作品第几次版权签名认证以及是否是转权等
	//椭圆加密签名就是两个
	Contentkeyrtext string//作品内容加密标签r,
	Contentkeystext string//作品内容加密标签s,
	Authorizationmoney int //版权授权费用
	Transactionmoney int //版权转权费用
	Transactionhash string
	Time time.Time//时间戳
}//主要用户记录作品的具体内容,不存储到区块链中,只在数据库文件中存储


type Works struct {
	Works map[string]*Work//作品是一个结构体
}

type Workcontents struct {
	Workcontents map[string]*Workcontent//作品是一个结构体
}

func NewWorks()  (*Works,error) {
	works := Works{}
	works.Works = make(map[string]*Work)
	//fmt.Println("加载文件")
	err := works.LoadworkFromFile()
	return &works,err
}//从文件中读取work信息进行操作,返回works是一个对于work的映射

func NewWorkcontents()  (*Workcontents,error) {
	workcontents := Workcontents{}
	workcontents.Workcontents = make(map[string]*Workcontent)
	//fmt.Println("加载文件")
	err := workcontents.LoadworkcontentFromFile()
	return &workcontents,err
}//从文件中读取workcontent信息进行操作,返回workcontents是一个对于workcontent的映射

func NewWork(Useraddress,Title,Abstract,Content,Contentrkey,Contentskey,nodeID string,Authorizationmoney,Transactionmoney int,transactionhash,simhash string) *Work  {
	wallets,err := NewWallets(nodeID)//创建一个钱包集合,并通过wallet.dat文件读取到已有的钱包集合
	if err != nil {
		log.Panic(err)
	}
	_wallet := wallets.GetWallet(Useraddress)//通过Useraddress得到用户信息
	work := Work{_wallet.Userinformation.Idcard,Useraddress,_wallet.Userinformation.Name, Title,
		Abstract,Content,Contentrkey,Contentskey,Authorizationmoney,Transactionmoney,transactionhash,time.Now(),simhash}
	return &work
}//创建一个带有用户地址的work


func SignECCWork(msg []byte, privKey ecdsa.PrivateKey)([]byte,[]byte) {//msg是要签名的内容
	//计算哈希值
	hash := sha256.New()
	//填入数据
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//对哈希值生成数字签名
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, bytes)
	if err != nil {
		panic(err)
	}
	rtext, _ := r.MarshalText()
	stext, _ := s.MarshalText()
	return rtext, stext
}//针对作品信息进行私钥签名

//验证数字签名 , msg是要签名的内容
func VerifySignECCWork(msg []byte,rtext,stext []byte,publicKey ecdsa.PublicKey) bool{
	//计算哈希值
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//验证数字签名
	var r,s big.Int
	r.UnmarshalText(rtext)
	s.UnmarshalText(stext)
	verify := ecdsa.Verify(&publicKey, bytes, &r, &s)
	return verify
}//通过公钥,签名标签以及签名的内容验证是否是正确的签名



func (ws *Works) CreateWork(Useraddress,Title,Abstract,content,Contentrkey,Contentskey,nodeID string,Authorizationmoney,Transactionmoney int,transactionhash,simhash string) string {
	work := NewWork(Useraddress,Title,Abstract,content,Contentrkey,Contentskey,nodeID,Authorizationmoney,Transactionmoney,transactionhash,simhash)
	biaoqian := fmt.Sprintf("%s",work.GetSHA256biaoqian([]byte(Useraddress+Title)))//获得每个作品对应的标签

	ws.Works[biaoqian] = work
	return biaoqian
}//创建作品work返回对于作品的键

func (w Work)GetSHA256biaoqian(message []byte)string{
	//方法一：
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode

}//生成每个作品对应的作品标签

func (w Workcontent)GetSHA256Signbiaoqian(message []byte)string{
	//方法一：
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode

}//生成每个作品对应的作品标签


func (ws *Works) LoadworkFromFile() error  {
	if _, err := os.Stat(workFile);os.IsNotExist(err) {//找到指定文件判断文件是否存在,
		// Stat返回描述指定文件的FileInfo(接口)。FileInfo描述一个文件，并由Stat和Lstat返回。
		// IsNotExist返回一个布尔值，指示是否知道该错误报告文件或目录不存在。它既满足了一些系统错误，也满足了一些错误。
		return err
	}

	fileContent, err := ioutil.ReadFile(workFile)//ReadFile读取按文件名命名的文件并以字节的形式返回内容。
	if err != nil {
		log.Panic(err)
	}

	var works Works
	gob.Register(elliptic.P256()) //P256返回一条实现P-256的曲线(参见FIPS 186-3, D.2.3节)
	// 加密操作是使用常量时间算法实现的。并且进行注册
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))//NewDecoder返回一个从io.Reader中读取数据的新解码器。
	// 解码器管理从连接的远程端读取的类型和数据信息的接收。
	err = decoder.Decode(&works)
	if err != nil {
		log.Panic(err)
	}
	ws.Works = works.Works
	return nil
}//从work.date文件中读取信息

func (ws Works) SaveworkToFile()  {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)//NewEncoder返回一个将在io.Writer上传输的新编码器。
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(workFile, content.Bytes(),0644)
	if err != nil {
		log.Panic(err)
	}
}//将作品版权注册信息保存到work.date文件


//作品内容的加载
func (ws *Workcontents) LoadworkcontentFromFile() error  {
	if _, err := os.Stat(workcontentFile);os.IsNotExist(err) {//找到指定文件判断文件是否存在,
		// Stat返回描述指定文件的FileInfo(接口)。FileInfo描述一个文件，并由Stat和Lstat返回。
		// IsNotExist返回一个布尔值，指示是否知道该错误报告文件或目录不存在。它既满足了一些系统错误，也满足了一些错误。
		return err
	}

	fileContent, err := ioutil.ReadFile(workcontentFile)//ReadFile读取按文件名命名的文件并以字节的形式返回内容。
	if err != nil {
		log.Panic(err)
	}

	var workcontents Workcontents
	gob.Register(elliptic.P256()) //P256返回一条实现P-256的曲线(参见FIPS 186-3, D.2.3节)
	// 加密操作是使用常量时间算法实现的。并且进行注册
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))//NewDecoder返回一个从io.Reader中读取数据的新解码器。
	// 解码器管理从连接的远程端读取的类型和数据信息的接收。
	err = decoder.Decode(&workcontents)
	if err != nil {
		log.Panic(err)
	}
	ws.Workcontents = workcontents.Workcontents
	return nil
}//从workcontent.date文件中读取信息

func (ws Workcontents) SaveworkcontentToFile()  {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)//NewEncoder返回一个将在io.Writer上传输的新编码器。
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(workcontentFile, content.Bytes(),0644)
	if err != nil {
		log.Panic(err)
	}
}//将作品版权注册以及作品内容信息保存到workcontent.date文件


////通过用户地址以及作品题目找到作品信息
//func (cli *CLI)Getworkbyaddressandtitle(address string,title ,nodeID string) Work  {
//
//	work := Getworkintransaction(address,title,nodeID)
//
//	return work
//}
//
////通过用户地址以及作品题目找到作品内容信息(认证时的作品记录)
//func (cli *CLI)Getworkcontentbyaddressandtitle(address string,title string) Workcontent  {
//	workcontents, _ := NewWorkcontents()//返回一个结构体,只包含一个string到works的映射,并且加载work文件中的作品信息,
//	workcontent := Workcontent{}//定义一个空的work用于调用GetSha256biaoqian函数获得作品的标签
//	biaoqian := workcontent.GetSHA256Signbiaoqian([]byte(address+title))//通过用户地址以及作品标题读取到作品的标签,以便找到作品规定的转权费用
//	if workcontents.Workcontents[biaoqian] == nil {
//		return Workcontent{}
//	}else {
//		return  *workcontents.Workcontents[biaoqian]
//	}
//}



// 得到存储在works里面的作品的键
func (wk *Works) GetWorkkeys() []string {
	var keys []string
	for key := range wk.Works {
		keys = append(keys, key)
	}
	return keys
}