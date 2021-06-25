package core

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12//命令长度

var nodeAddress string//挖矿的节点窗口(3002)
var miningAddress string//挖矿的地址
var knownNodes = []string{"localhost:3000"}//节点的地址
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)//内存池,用来存储交易
var mining = false

type addr struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

//给我看一下你有什么区块,不是全部区块,而是请求了一个块哈希的列表(减轻网络负担)
type getblocks struct {
	AddrFrom string
}

//用于某个块或交易的请求，它可以仅包含一个块或交易的 ID
type getdata struct {
	AddrFrom string
	Type     string//表明类型,展示的是块还是交易
	ID       []byte
}

//使用inv向其他节点展示当前节点有什么块和交易,仅仅是哈希值,不是完整的区块链和交易
type inv struct {
	AddrFrom string
	Type     string//表明类型,展示的是块还是交易
	Items    [][]byte
}

type tx struct {
	AddFrom     string
	Transaction []byte
}


//新节点可以给其他节点发送消息
type verzion struct {
	Version    int//此处只有一个版本  用于找到一个更长的区块链
	BestHeight int64//节点高度,存储的区块数量
	AddrFrom   string//发送节点的地址,并不是用户地址
}


//消息在底层使用字节序列来实现的,命令变成字节
func commandToBytes(command string) []byte {
	var bytes [commandLength]byte
	//创建一个12字节的缓冲区,并用命令名来进行填充,将剩下的字节置为空

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

//从字节序列提取出命令
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}//bytesToCommand 用来提取命令名，并选择正确的处理器处理命令主体：

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func sendAddr(address string) {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *Block) {
	fmt.Println("sendBlock输出如下:","addr:",addr)
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func sendData(addr string, data []byte) {//addr节点地址,3000
	fmt.Println("sendData输出如下:","addr为",addr)
    conn, err := net.Dial(protocol, addr)//dial拨号连接到指定网络上的地址。
	//Conn是一种通用的面向流的网络连接。是一个接口
	//多个goroutine可以同时调用Conn上的方法。
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))//将副本从data复制到conn，直到在data上达到EOF或发生错误。
	// 它返回复制的字节数和复制时遇到的第一个错误（如果有）。
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {//items是一个块的哈希列表或者是交易的列表
	fmt.Println("sendInv输出如下:","addr:",address)
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

//向version发生节点发送更新区块请求
func sendGetBlocks(address string) {//address是version发送节点的地址
	fmt.Println("sendGetBlocks输出如下:","addr:",address)
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	fmt.Println("sendGetData输出如下:","addr:",address)
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(addr string, tnx *Transaction) {
	fmt.Println("sendTx输出如下:","addr:",addr)
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)//序列化data为字节数组
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}


//发送version消息来查询是否过期
func sendVersion(addr string, bc *Blockchain) {//
    fmt.Println("sendVersion输出如下:","addr如下:",addr)
	bestHeight := bc.GetBestheight()//发送节点的区块高度

	payload := gobEncode(verzion{nodeVersion, bestHeight, nodeAddress})//后面的字节包含了god编码的信息结构

	request := append(commandToBytes("version"), payload...)//前十二个字节指定了命令名,将Version消息序列化后添加上命令名

	sendData(addr, request)//addr是接收节点信息
}

func handleAddr(request []byte) {
	fmt.Println("handleAddr输出如下:")
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}

func handleBlock(request []byte, bc *Blockchain) {
	fmt.Println("handleBlock输出如下:")
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlcok(blockData)//区块数据反序列化为一个区块

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)//当接收到一个新块时，我们把它放到区块链里面

	fmt.Printf("Added block %x\n", block.Hash)

	//添加新区快时的时间
	fmt.Println("中心节点同步区块之后的时间")
	t2 := time.Now().Unix()
	fmt.Println(t2)
	//时间戳转化为具体时间
	fmt.Println(time.Unix(t2, 0).String())

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:] //如果还有更多的区块需要下载，我们继续从上一个下载的块的那个节点继续请求
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()//当最后把所有块都下载完后，对 UTXO 集进行重新索引
	}
}

//处理INV命令的函数
func handleInv(request []byte, bc *Blockchain) {
	fmt.Println("handleInv输出如下:")
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items//如果收到了块哈希，保存在这个变量中（二维数组的形式），来跟踪已下载的块。这是为了从不同节点下载块。


		blockHash := payload.Items[0]//
		sendGetData(payload.AddrFrom, "block", blockHash)//在将块置于传送状态时，我们给 inv 消息的发送者发送 getdata 命令并更新 blocksInTransit

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {//把blockHash之外的b全部导入到newINTransit中
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]//只拿到第一个哈希，因为这里不会发送有多重哈希的inv。
//
		if mempool[hex.EncodeToString(txID)].ID == nil {//如果内存池中没有这个哈希，就发送getdata消息
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

//处理getblocks命令
func handleGetBlocks(request []byte, bc *Blockchain) {
	fmt.Println("handleGetBlocks输出如下:")
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()//返回链中所有快的hash列表
	sendInv(payload.AddrFrom, "block", blocks)//返回所有块的哈希列表
}

//GetData的处理器
func handleGetData(request []byte, bc *Blockchain) {
	fmt.Println("handleGetData输出如下:")
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])//缓冲区
	dec := gob.NewDecoder(&buff)//解码器
	err := dec.Decode(&payload)//解码
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {//如果请求块，则返回一个块
		block, err := bc.GetBlock([]byte(payload.ID))//通过区块hash找到区块
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)//将区块发送给消息发起者
	}

	if payload.Type == "tx" {//如果它们请求一笔交易，则返回交易
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
	}
}

func handleTx(request []byte, bc *Blockchain) {
	fmt.Println("handleTx输出如下:")
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	//表明tx的来源
	fmt.Println("tx来自于",payload.AddFrom)
	txData := payload.Transaction
	tx := DeserializeTransaction(txData)//反序列化得到一个交易实体
	mempool[hex.EncodeToString(tx.ID)] = tx//将新交易放到内存池中,

	if nodeAddress == knownNodes[0] {//检查当前节点是否是中心节点
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID}) //将新的交易推送给网络中的其他节点
			}
		}
	} else {
		if len(mempool) >= 1 && len(miningAddress) > 0 && mining == false{//矿工代码块.
			//如果当前节点（矿工）的内存池中有两笔或更多的交易，开始挖矿
		MineTransactions:
			var txs []*Transaction

			for id := range mempool {
				tx := mempool[id]
				//if bc.VerifyTransaction(&tx) {//内存池中所有交易都是通过验证的
				//	txs = append(txs, &tx)
				//	fmt.Println("有一个交易验证成功!!")
				//}
				txs = append(txs,&tx)
			}//内存池中的交易依次转到txs中

			if len(txs) == 0 {//如果没有有效交易，则挖矿中断
				fmt.Println("所有交易无效！ 等待新的交易...")
				return
			}



			//挖矿奖励
			cbTx := NewCoinbaseTX(miningAddress, "挖矿奖励")//附带奖励的 coinbase 交易
			txs = append(txs, cbTx)//验证后的交易被放到一个块里

			UTXOSet := UTXOSet{bc}

			////挖矿之前的时间
			//fmt.Println("挖矿之前的时间:")
			//t1 := time.Now().Unix()
			//fmt.Println(t1)
			////时间戳转化为具体时间
			//fmt.Println(time.Unix(t1, 0).String())

			mining = true
			newBlock := bc.MineBlock(txs)// MineBlock使用提供的transaction事务挖掘新的区块
			UTXOSet.Update(newBlock)//当块被挖出来以后，UTXO 集会被重新索引。
			mining = false

			fmt.Println("New block is mined!区块高度:",newBlock.Height,"区块的交易数:",newBlock.numbertx())

			//for _, tx := range txs {
			//	txID := hex.EncodeToString(tx.ID)
			//	delete(mempool, txID)//当一笔交易被挖出来以后，就会被从内存池中移除
			//}

			for _, tx := range txs {
				for _,t := range newBlock.Transactions {
					if tx == t {
						txID := hex.EncodeToString(tx.ID)
						delete(mempool, txID)//当一笔交易被挖出来以后，就会被从内存池中移除
						break
					}
				}
			}//当交易池中的交易在新的区块中的时候再删除交易池的交易

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})//当前节点所连接到的所有其他节点，接收带有新块哈希的 inv 消息
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}//矿工代码块
	}
}

//Version 命令处理器
func handleVersion(request []byte, bc *Blockchain) {
	fmt.Println("handleVersion输出如下:")
	var buff bytes.Buffer//Buffer是一个可变大小的字节缓冲区，具有Read和Write方法。
	var payload verzion

	buff.Write(request[commandLength:])// Write将p的内容附加到缓冲区，根据需要增长缓冲区。
	dec := gob.NewDecoder(&buff)//NewDecoder返回一个从io.Reader读取的新解码器。
	err := dec.Decode(&payload)//对请求进行解码,提取有效的信息.解码从输入流中读取下一个值，并将其存储在由空接口值表示的数据中。
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestheight()
	foreignerBestHeight := payload.BestHeight//从节点中提取区块高度信息

	//与自身区块高度进行对比
	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)//如果自身的区块高度小,则更新区块
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)//自身区块高度大,则回复version消息
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}
//此函数用来处理节点接收到的命令
func handleConnection(conn net.Conn, bc *Blockchain) {
	//Conn是一种通用的面向流的网络连接(一个接口)。 多个goroutine可以同时调用Conn上的方法。

	request, err := ioutil.ReadAll(conn)//ReadAll从r读取，直到出现错误或EOF并返回它读取的数据。成功的调用返回err == nil
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])//从字节序列中提取出命令
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

// StartServer 启动一个节点
func StartServer(nodeID, minerAddress string) {//minerAddress指定了接受挖矿奖励的地址
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress//miningAddress是一个全局变量
	ln, err := net.Listen(protocol, nodeAddress)//TCP监听,ln监听器是面向流的协议的通用网络监听器。
	//多个goroutine可以同时调用Listener上的方法。
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockchain(nodeID)//得到对应节点的本地区块链

	//如果当前节点不是中心节点,必须向中心节点发送version消息来查询是否自己的区块链已经过时
	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)//向中心节点(3000)发送Version消息
	}

	for {
		conn, err := ln.Accept()// Accept等待并返回与侦听器的下一个连接。
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func gobEncode(data interface{}) []byte {//将data接口序列化为字节数组
	var buff bytes.Buffer//Buffer缓冲区是具有读和写方法的可变大小的字节缓冲区。

	enc := gob.NewEncoder(&buff)//创建基于buffer内存的编码器
	err := enc.Encode(data)//序列化,使用编码器对data进行编码
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()//返回buff的字节数组
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}
