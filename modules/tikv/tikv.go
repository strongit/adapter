package tikv

import (
	"adapter/modules/conf"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
	"log"
)

type kv struct {
	Key   string
	Value string
}

var Client *tikv.RawKVClient

// InitStore 初始化
func Init() {
	var err error
	Client, err = tikv.NewRawKVClient([]string{conf.RunTimeInfo.PDHost}, config.Security{})
	if err != nil {
		log.Println(err)
	}
}

// Puts 写入数据
func Puts(args ...[]byte) error {
	for i := 0; i < len(args); i += 2 {
		key, val := args[i], args[i+1]
		err := Client.Put(key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

//BatchPut 批量写入
func BatchPut(keys, values [][]byte) error {
	err := Client.BatchPut(keys, values)
	if err != nil {
		return err
	}
	return nil
}


// Dels 删除数据
func Dels(keys ...[]byte) error {
	for i := 0; i < len(keys); i += 1 {
		err := Client.Delete(keys[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Delall 批量删除
func Delall(startKey []byte, endKey []byte, limit int) error {
	keys, _, err := Client.Scan(startKey, limit)
	if err != nil {
		return err
	}
	for i := 0; i < len(keys); i += 1 {
		Dels(keys[i])
	}
	return nil
}

// Get 获取数据
func Get(k []byte) (kv, error) {
	v, err := Client.Get(k)
	if err != nil {
		return kv{}, err
	}
	return kv{Key: string(k), Value: string(v)}, nil
}

// Scan 批量获取数据
func Scan(startKey []byte, endKey []byte, limit int) ([]kv, error) {
	var kvs []kv
	keys, values, err := Client.Scan(startKey, limit)
	if err != nil {
		return kvs, err
	}
	for i := 0; i < len(keys); i += 1 {
		kvs = append(kvs, kv{Key: string(keys[i]), Value: string(values[i])})
	}
	return kvs, nil
}
