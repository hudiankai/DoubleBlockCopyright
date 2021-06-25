package core

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const zhuanquanFile = "workzhuanquan.data"//存储写入到区块链中的作品的版权转权


type Zhuanquan struct {
	Buyaddress string
	Selladdress string
	WorkTital string
	ZhuanquanMoney int
	Time time.Time
	Newauthorizationmoney int
	Newtransactionmoney int
	Transactionhash string
	Remarks string//上一任版权所有者的版权记录所在的区块hash
}//作品转权信息


type Zhuanquans struct {
	Zhuanquans map[string]*Zhuanquan//授权记录是一个结构体
}

func NewZhuanquans()  (*Zhuanquans,error) {
	zhuanquans := Zhuanquans{}
	zhuanquans.Zhuanquans = make(map[string]*Zhuanquan)
	err := zhuanquans.LoadzhuanquanFromFile()
	return &zhuanquans,err
}//从文件中读取work信息进行操作,返回works是一个对于work的映射


func (s Zhuanquan)GetSHA256Zhuanquanbiaoqian(message []byte)string{
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

}//生成每个作品授权记录对应的标签


func (s *Zhuanquans) LoadzhuanquanFromFile() error  {
	if _, err := os.Stat(zhuanquanFile);os.IsNotExist(err) {//找到指定文件判断文件是否存在,
		// Stat返回描述指定文件的FileInfo(接口)。FileInfo描述一个文件，并由Stat和Lstat返回。
		// IsNotExist返回一个布尔值，指示是否知道该错误报告文件或目录不存在。它既满足了一些系统错误，也满足了一些错误。
		return err
	}

	fileContent, err := ioutil.ReadFile(zhuanquanFile)//ReadFile读取按文件名命名的文件并以字节的形式返回内容。
	if err != nil {
		log.Panic(err)
	}

	var zhuanquans Zhuanquans
	gob.Register(elliptic.P256()) //P256返回一条实现P-256的曲线(参见FIPS 186-3, D.2.3节)
	// 加密操作是使用常量时间算法实现的。并且进行注册
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))//NewDecoder返回一个从io.Reader中读取数据的新解码器。
	// 解码器管理从连接的远程端读取的类型和数据信息的接收。
	err = decoder.Decode(&zhuanquans)
	if err != nil {
		log.Panic(err)
	}
	s.Zhuanquans = zhuanquans.Zhuanquans

	return nil
}//从shouquan.data文件中读取信息

func (ws *Zhuanquans) SavezhuanquanToFile()  {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)//NewEncoder返回一个将在io.Writer上传输的新编码器。
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(zhuanquanFile, content.Bytes(),0644)
	if err != nil {
		log.Panic(err)
	}
}//将作品版权注册信息保存到shoquan.data文件

