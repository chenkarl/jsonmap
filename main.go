package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	filepath := "./data.json"
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	var ob interface{}
	err = json.Unmarshal(data, &ob)
	if err != nil {
		return
	}
	es, er := SetType(ob)
	if er != nil {
		return
	}
	esData, err := json.Marshal(es)
	if err != nil {
		return
	}
	ioutil.WriteFile("./target.json", esData, 0644)
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
			obmap, ok := value.([]interface{})[0].(map[string]interface{})
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
