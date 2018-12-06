package gredis

import (
	"gredis/consts"
	"gredis/sds"
	"strconv"
	"strings"
	"time"
)

// 预先配置RedisSharedIntegers个共享整数编码的对象
var SharedIntObj []*Object

func init() {
	for i := 0; i < consts.RedisSharedIntegers; i++ {
		intObj := NewObject(consts.TypeRedisString, i)
		intObj.Encoding = consts.EncodingInt
		SharedIntObj = append(SharedIntObj, intObj)
	}
}

type Object struct {
	Type     uint8
	Encoding uint8
	Ptr      interface{}
	Lru      int64
}

func NewObject(typ uint8, ptr interface{}) *Object {
	obj := new(Object)
	obj.Type = typ
	obj.Ptr = ptr
	obj.Encoding = consts.EncodingRaw
	obj.Lru = time.Now().Unix()
	return obj
}

func NewRawStringObject(str string) *Object {
	return NewObject(consts.TypeRedisString, sds.NewSds(str))
}

func NewObjectFromLongValue(val int64) *Object {
	if val >= 0 && val < consts.RedisSharedIntegers {
		return SharedIntObj[val]
	}

	obj := NewObject(consts.TypeRedisString, val)
	obj.Encoding = consts.EncodingInt

	return obj
}

func CompareStringObject(obj1 *Object, obj2 *Object) int {
	if obj1.Type != consts.TypeRedisString || obj2.Type != consts.TypeRedisString {
		panic(consts.ErrCompareType)
	}

	var s1 string
	var s2 string

	switch v := obj1.Ptr.(type) {
	case *sds.Sds:
		s1 = v.String()
	case int:
		s1 = strconv.Itoa(v)
	}

	switch v := obj2.Ptr.(type) {
	case *sds.Sds:
		s2 = v.String()
	case int:
		s2 = strconv.Itoa(v)
	}

	return strings.Compare(s1, s2)
}

// 优化了CompareStringObject，可以不需要转化为string
func EqualStringObject(obj1 *Object, obj2 *Object) bool {
	if obj1.Encoding == consts.EncodingInt &&
		obj2.Encoding == consts.EncodingInt {
		return obj1.Ptr == obj2.Ptr
	}

	return CompareStringObject(obj1, obj2) == 0
}
