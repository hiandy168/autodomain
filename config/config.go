package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type JsonStruct struct{}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (self *JsonStruct) LoadJson(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(errors.New("读取配置文件错误!"))
		return
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, v)
	if err != nil {
		panic(errors.New("解析配置文件错误!"))
		return
	}
}

type Configdata struct {
	Time      int
	Id        string
	Secret    string
	Domain    string
	CheckUrls []string
}
