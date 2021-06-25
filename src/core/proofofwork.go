package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt32
)

const targetBits  = 16  //控制算法的难度，数越大难度越大,10表示哈希的前10位必须是0

type ProofOfWork struct {
	block *Block//当前要验证的区块
	target *big.Int//计算区块满足的目标，最终找到的哈希必须要小于目标
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target :=big.NewInt(1)//设置初始值为1的target
	target.Lsh(target,uint(256 - targetBits))//Lsh是移位操作，左移,target左移256-targetBits位再赋值给target

	pow := &ProofOfWork{b, target}

	return pow
}// 新建powwork，并且返回一个pow

//用来准备POW证明所需的数据,也可以用来验证工作量
func (pow *ProofOfWork) prepareData(nonce int )[]byte{
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),//IntToHex将int64类型的数据转换成字节切片
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
		)

	return data
}//拼接区块属性，返回字节数组

func (pow *ProofOfWork)Run() (int, []byte) {
	var hashInt big.Int//存储新生成hash值
	var hash [32]byte //存储hash值
	nonce := 0//初始化随机数(计数器)
	//fmt.Println(pow.block)

	//fmt.Printf("Mining the block comtaining \"%s\"\n",pow.block.Data)
	for nonce < maxNonce {//防止溢出的"无限"循环
		//将block的属性拼接成字节切片并返回,准备数据
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)//对数据进行hash计算,将拼接后的字节数组生成hash

		hashInt.SetBytes(hash[:])//将hash转换为Int类型(大整数),用hashInt存储,hash[:]是切片,hash是数组

		if hashInt.Cmp(pow.target) == -1 {//hashInt<taeget的话 就结束,(将大整数与目标进行比较)
			break
		}else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}//运行项目POW,寻找合适的nonce,并返回目标区块的有效hash


//验证工作量,只要哈希小于目标就是有效工作量,采用计数器的值直接与区块中的数据进行计算hash值
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])//将hash设置为无符号整形的字节,并赋值给hashInt,并返回hashInt

	isValid := hashInt.Cmp(pow.target) == -1//比较难度值,如果等于-1,则表示hashInt小于target

	return isValid
}
