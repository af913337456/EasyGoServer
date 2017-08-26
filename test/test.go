package test

import (
	"os"
	"fmt"
)

func checkF(fileName string) bool {
	_,err := os.Stat(fileName)
	if err!=nil {
		os.Create(fileName)
		return false
	}else{
		fmt.Println("已存在 "+fileName)
		return true
	}
}

func main()  {
	file,err := os.Create("struct.go")
	defer file.Close()
	if err!=nil {
		fmt.Println("outputStructFile 293")
		fmt.Println(err)
		return
	}
	file.WriteString("package main\n")
	file.Sync()


}
