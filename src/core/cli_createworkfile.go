package core

import (
	"fmt"
	"log"
	"os"
)

func (cli *CLI)CreateworkFile(filename ,nodeid string)  {
	f,err := os.Create(filename)
	defer f.Close()
	//
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(filename,"文件创建成功,可以进行创作作品")
	//else {
	//	_, err = f.Write([]byte(content))
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//}
}
