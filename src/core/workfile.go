package core

import (
"fmt"
"io/ioutil"
"log"
"os"
)

const work  = "workfile.data"

func ReadFile(filename string) string {

	f ,err := os.OpenFile(filename,os.O_RDONLY,0600)

	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}else {
		contentByte,error := ioutil.ReadAll(f)
		if error != nil{
			log.Panic(error)
		}
		return string(contentByte)
	}

	return "作品文件读取失败!"
}

func CreateFile(filename , content string)  {
	f,err := os.Create(filename)
	defer f.Close()

	if err != nil {
		log.Panic(err)
	}else {
		_, err = f.Write([]byte(content))
		if err != nil {
			log.Panic(err)
		}
	}
}

func Writefile ()  {

}

