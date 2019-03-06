package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func main() {

	defer func() {
		fatal := recover()
		if fatal != nil {
			fmt.Println("程序崩溃，错误信息为:", fatal)
			fmt.Println("按回车退出")
			fmt.Scanln()
		}
	}()
	filepath := "./data.json"
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("读文件错误")
		fmt.Println("按回车退出")
		fmt.Scanln()
		return
	}
	var ob interface{}
	err = json.Unmarshal(data, &ob)
	if err != nil {
		fmt.Println("json反格式化失败")
		fmt.Println("按回车退出")
		fmt.Scanln()
		return
	}
	es, er := SetType(ob)
	if er != nil {
		fmt.Println("格式转义错误")
		fmt.Println("按回车退出")
		fmt.Scanln()
		return
	}
	esData, err := json.Marshal(es)
	if err != nil {
		fmt.Println("json格式化失败")
		fmt.Println("按回车退出")
		fmt.Scanln()
		return
	}
	ioutil.WriteFile("./target.json", esData, 0644)
	fmt.Println("按回车退出")
	fmt.Scanln()
}

//SetType ...
func SetType(value interface{}) (es *EsValue, err interface{}) {
	log.Println(value)
	err = nil
	switch value.(type) {
	case int:
		{
			return NewEsLong(), err
		}
	case float64:
		{
			return NewEsLong(), err
		}
	case float32:
		{
			return NewEsLong(), err
		}
	case time.Time:
		{
			return NewEsDate(), err
		}
	case string:
		{
			return NewEsText(), err
		}
	case bool:
		{
			return NewEsBoolean(), err
		}
	case map[string]interface{}:
		{
			obmap, ok := value.(map[string]interface{})
			kv := make(EsValue)
			if ok {
				for k, v := range obmap {
					esMap, err := SetType(v)
					if err != nil {
						log.Fatalln(err)
						return NewEsError(), err
					}
					log.Println(k, esMap)
					kv[k] = esMap
				}
			}
			return &kv, err
		}
	case []interface{}:
		{
			arr := value.([]interface{})
			if len(arr) == 0 {
				// err = "该数组无数据"
				return NewEsError(), err
			}
			obmap, ok := arr[0].(map[string]interface{})
			kv := make(EsValue)
			if ok {
				for k, v := range obmap {
					esMap, err := SetType(v)
					if err != nil {
						log.Fatalln(err)
						return NewEsError(), err
					}
					kv[k] = esMap
				}
			}
			return &kv, err
		}
	default:
		{
			err = "无该数据类型"
			return NewEsError(), err
		}
	}
}

//EsValue ... EsValue
type EsValue map[string]interface{}

// NewEsError ... 获取一个EsError
func NewEsError() *EsValue {
	ev := make(EsValue)
	ev["type"] = "error"
	return &ev
}

// NewEsLong ... 获取一个EsLong
func NewEsLong() *EsValue {
	ev := make(EsValue)
	ev["type"] = "long"
	return &ev
}

// NewEsBoolean ... 获取一个EsBoolean
func NewEsBoolean() *EsValue {
	ev := make(EsValue)
	ev["type"] = "boolean"
	return &ev
}

// NewEsDate ... 获取一个ESDate
func NewEsDate() *EsValue {
	ev := make(EsValue)
	ev["type"] = "date"
	return &ev
}

// NewEsText ... 获取一个EsText
func NewEsText() *EsValue {
	ev := make(EsValue)
	fields := make(map[string]interface{})
	keyword := make(map[string]interface{})
	ev["type"] = "text"
	keyword["ignore_above"] = 256
	keyword["type"] = "keyword"
	fields["keyword"] = keyword
	ev["fields"] = fields
	return &ev
}
