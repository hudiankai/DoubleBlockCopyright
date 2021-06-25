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

const shouquanFile = "workshouquan.date"//存储写入到区块链中的作品的版权授权记录


type Shouquan struct {
	Buyaddress string
	Selladdress string
	Title string
	Money int
	Remarks string//存储版权注册的交易所在的区块hash
	Transactionhash string
	Time time.Time
	Workcontent string
}//作品授权信息


type Shouquans struct {
	Shouquans map[string]*Shouquan//授权记录是一个结构体
}

func NewShouquans()  (*Shouquans,error) {
	shouquans := Shouquans{}
	shouquans.Shouquans = make(map[string]*Shouquan)
	err := shouquans.LoadshouquanFromFile()
	return &shouquans,err
}//从文件中读取work信息进行操作,返回works是一个对于work的映射


func (s Shouquan)GetSHA256Shouquanbiaoqian(message []byte)string{
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


func (s *Shouquans) LoadshouquanFromFile() error  {
	if _, err := os.Stat(shouquanFile);os.IsNotExist(err) {//找到指定文件判断文件是否存在,
		// Stat返回描述指定文件的FileInfo(接口)。FileInfo描述一个文件，并由Stat和Lstat返回。
		// IsNotExist返回一个布尔值，指示是否知道该错误报告文件或目录不存在。它既满足了一些系统错误，也满足了一些错误。
		return err
	}

	fileContent, err := ioutil.ReadFile(shouquanFile)//ReadFile读取按文件名命名的文件并以字节的形式返回内容。
	if err != nil {
		log.Panic(err)
	}

	var shouquans Shouquans
	gob.Register(elliptic.P256()) //P256返回一条实现P-256的曲线(参见FIPS 186-3, D.2.3节)
	// 加密操作是使用常量时间算法实现的。并且进行注册
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))//NewDecoder返回一个从io.Reader中读取数据的新解码器。
	// 解码器管理从连接的远程端读取的类型和数据信息的接收。
	err = decoder.Decode(&shouquans)
	if err != nil {
		log.Panic(err)
	}
	s.Shouquans = shouquans.Shouquans

	return nil
}//从shouquan.data文件中读取信息

func (ws *Shouquans) SaveshouquanToFile()  {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)//NewEncoder返回一个将在io.Writer上传输的新编码器。
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(shouquanFile, content.Bytes(),0644)
	if err != nil {
		log.Panic(err)
	}
}//将作品版权注册信息保存到shoquan.data文件


