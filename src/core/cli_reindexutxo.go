package core

import (
	"fmt"
)

//查找UTXO集中的交易数
func (cli *CLI) reindexUTXO(nodeID string) {
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex() //在现实中如果能保证自己下载的链节点是完整的，可以忽略。

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done!!! There are %d transactions in the UTXO set.\n", count)
}

