package dict

import (
	"sync"
	"errors"
	"math/rand"
	"time"
)

var ErrDictAddErr = errors.New("key is exist")

type Dict struct {
	sync.Map
}

func NewDict() *Dict {
	return new(Dict)
}

func (dict *Dict) Add(key, val interface{}) error {
	_, ok := dict.LoadOrStore(key, val)
	if ok {
		return ErrDictAddErr
	}
	return nil
}

func (dict *Dict) Replace(key, val interface{}) {
	dict.Store(key, val)
}

func (dict *Dict) FetchKey(key interface{}) interface{} {
	result, ok := dict.Load(key)
	if ok {
		return result
	}
	return nil
}

func (dict *Dict) GetRandomKey() map[interface{}]interface{} {
	var keys []interface{}

	dict.Range(func(key, _ interface{}) bool {
		keys = append(keys, key)
		return true
	})

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(keys))

	key := keys[index]

	val, ok := dict.Load(key)
	if !ok {
		return nil
	}

	return map[interface{}]interface{}{
		key: val,
	}
}


