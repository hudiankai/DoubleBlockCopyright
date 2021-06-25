package core

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet//钱包是一个结构体,包含一个公钥和一个私钥
}

func NewWallets(nodeID string)  (*Wallets,error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile(nodeID)

	return &wallets,err
}


func (ws *Wallets) CreateWalletuser(Name, Role, Idcard, password string ) string {
	wallet := NewWalletuser(Name, Role, Idcard,password)
	address := fmt.Sprintf("%s",wallet.GetAddress())
	ws.Wallets[address] = wallet

	fmt.Println("用户钱包创建成功:")
	fmt.Println("用户姓名:",wallet.Userinformation.Name)
	fmt.Println("用户备注:",wallet.Userinformation.Role)
	fmt.Println("用户ID:",wallet.Userinformation.Idcard)
	fmt.Println("用户创建时间:",wallet.Userinformation.Time)

	return address
}//创建带有用户信息的钱包


func (ws *Wallets) CreateWallet( ) string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s",wallet.GetAddress())

	ws.Wallets[address] = wallet
	return address
}//创建不带用户信息的钱包


func (ws Wallets) GetWallet(address string) Wallet  {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile(nodeID string) error  {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile);os.IsNotExist(err) {//找到指定文件判断文件是否存在,
		// Stat返回描述指定文件的FileInfo(接口)。FileInfo描述一个文件，并由Stat和Lstat返回。
		// IsNotExist返回一个布尔值，指示是否知道该错误报告文件或目录不存在。它既满足了一些系统错误，也满足了一些错误。
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)//ReadFile读取按文件名命名的文件并以字节的形式返回内容。
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256()) //P256返回一条实现P-256的曲线(参见FIPS 186-3, D.2.3节)
	// 加密操作是使用常量时间算法实现的。并且进行注册
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))//NewDecoder返回一个从io.Reader中读取数据的新解码器。
	// 解码器管理从连接的远程端读取的类型和数据信息的接收。
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}

func (ws Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(walletFile, nodeID)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)//NewEncoder返回一个将在io.Writer上传输的新编码器。
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// 得到存储在wallets里的地址
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func (cli *CLI)Getaddressbypassword(password,name,nodeID string) string {
	wallets, _ := NewWallets(nodeID)
	for k,v :=range	wallets.Wallets {
		if v.Userinformation.Password==password && v.Userinformation.Name == name{
			//fmt.Println(k)
			return k
		}
	}
	return "该身份证未注册,没有对于的钱包地址"
}//通过身份证号查找用户对应的钱包地址


func Getwalletbyaddress(address,nodeID string) Wallet {
	wallets, _ := NewWallets(nodeID)
	//fmt.Println(wallets.Wallets[address])
	return *wallets.Wallets[address]
}//通过身份证号查找用户对应的钱包