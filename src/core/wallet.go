package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
	"time"
)

const version  = byte(0x00)

const addressChesksumLen  = 4


type Wallet struct {
	PrivateKey ecdsa.PrivateKey//PrivateKey表示ECDSA私钥,是一个结构体
	PublicKey []byte
	Userinformation User
}

type User struct {
	Name string
	Role string //用户角色(创作者,购买者,查询者)
	Idcard string
	Password string
	Time time.Time//用户注册时间
}//用户信息


func NewWalletuser(Name, Role, Idcard,password string) *Wallet  {
	private, public := newKeyPair()//创建公钥私钥对	var user = User{Name,Role,UserID,Idcard,time.Now()}
	var user = User{Name,Role,Idcard,password ,time.Now()}
	wallet := Wallet{private ,public,user}
	return &wallet
}//创建带有用户信息的钱包


func NewWallet() *Wallet  {
	private, public := newKeyPair()//创建公钥私钥对	var user = User{Name,Role,UserID,Idcard,time.Now()}
	var user = User{}
	wallet := Wallet{private ,public,user}
	return &wallet
}//创建不带用户信息的钱包


//生成一个地址(base58)
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionePayload := append([]byte{version},pubKeyHash...)//给哈希加上(地址生成算法的)版本前缀
	checksum := checksum(versionePayload)//使用SHA256对加了前缀的哈希进行再哈希,计算校验和,这是最终结果的后四个字节

	fullPayload := append(versionePayload, checksum...)//将校验和附加到version+PubKeyhash的组合中
	address := Base58Encode(fullPayload)//使用Base58对version+Pubkeyhash+checksum组合进行编码

	return address
}


func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)//Sum256返回数据的SHA256校验和(32位的自己数组)。

	RIPEMD160Hasher := ripemd160.New()//New返回一个计算校验和的新hash。
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])//对公钥哈希两次
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)//Sum将当前hash附加到b并返回结果字节切片。

	return publicRIPEMD160
}

//基于椭圆曲线生成一个新的秘钥对
func newKeyPair() (ecdsa.PrivateKey, []byte)  {
	curve := elliptic.P256()//需要一个椭圆曲线

	private, err := ecdsa.GenerateKey(curve,rand.Reader)//加密函数,GenerateKey是go加密库里面的函数，使用椭圆曲线生成一个私钥
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)//从私钥生成公钥,公钥是曲线上的点(x,y坐标的组合)

	return *private,pubKey
}

//校验位checksum,双重哈希运算
func checksum(payload []byte) []byte {
	//下面双重哈希payload，在调用中，所引用的payload为（version + Pub Key Hash）
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])//使用SHA256对加了前缀的哈希进行再哈希

	//addressChecksumLen代表保留校验位长度
	return secondSHA[:addressChesksumLen]
}

//判断输入的地址是否有效,主要是检查后面的校验位是否正确
func ValidateAddress(address string) bool {
	//解码base58编码过的地址
	pubKeyHash := Base58Decode([]byte(address))
	//拆分pubKeyHash,pubKeyHash组成形式为：(一个字节的version) + (Public key hash) + (Checksum)
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChesksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-addressChesksumLen]
	targetChecksum := checksum(append([]byte{version},pubKeyHash...))
	//比较拆分出的校验位与计算出的目标校验位是否相等
	return bytes.Compare(actualChecksum,targetChecksum) == 0
}