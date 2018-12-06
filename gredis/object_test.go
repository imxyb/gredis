package gredis

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"gredis/consts"
	"github.com/stretchr/testify/assert"
	"gredis/sds"
)

type ObjectTestSuite struct {
	suite.Suite
	Object *Object
}

func TestNewRawStringObject(t *testing.T) {
	obj := NewRawStringObject("abc")
	assert.EqualValues(t, consts.EncodingRaw, obj.Encoding)
	assert.EqualValues(t, "abc", obj.Ptr.(*sds.Sds).String())
	assert.EqualValues(t, consts.TypeRedisString, obj.Type)
}

func TestNewObjectFromLongValue(t *testing.T) {
	obj1 := NewObjectFromLongValue(1)
	assert.Equal(t, SharedIntObj[1], obj1)

	var pass bool
	obj2 := NewObjectFromLongValue(consts.RedisSharedIntegers + 100)
	for _, item := range SharedIntObj {
		if item == obj2 {
			pass = true
			break
		}
	}

	assert.False(t, pass)
	assert.EqualValues(t, consts.TypeRedisString, obj2.Type)
	assert.EqualValues(t, consts.RedisSharedIntegers + 100, obj2.Ptr)
	assert.EqualValues(t, consts.EncodingInt, obj2.Encoding)
}

func TestLongValueObjectCompare(t *testing.T) {
	obj1 := NewObjectFromLongValue(200000)
	obj2 := NewObjectFromLongValue(200000)
	assert.Equal(t, obj1, obj2)

	obj1 = NewObjectFromLongValue(2)
	obj2 = NewObjectFromLongValue(3)
	num1 := obj1.Ptr.(int)
	num2 := obj2.Ptr.(int)
	assert.True(t, num1 < num2)
}