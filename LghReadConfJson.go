
package main

/**

仅做读取 配置 json 文件作用

*/

import (
	io "io/ioutil"
	json "encoding/json"

)

type JsonStruct struct{

}

func ConfigJsonStruct () *JsonStruct {
	return &JsonStruct{}
}

func (self *JsonStruct) Load (filename string, v interface{}) {

	data, err := io.ReadFile(filename)

	if err != nil {
		return
	}

	datajson := []byte(data)

	err = json.Unmarshal(datajson, v)

	if err != nil{
		return
	}

}

type ServerConfig struct{

	Host string
	Port string

	FilePort string

	DbName string
	DbUser string
	DbPw   string
	DbPort string
}


var serverConfig ServerConfig

func bindServerConfig() {
	JsonParse := ConfigJsonStruct()
	/** 传入的 结构体 要和 json 的格式对上,否则返回是 null */
	JsonParse.Load("conf.json", &serverConfig)

}
