package zskiplist

import (
	"gredis/consts"
	"gredis/gredis"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ZskipListTestSuite struct {
	suite.Suite
	list *List
}

func (suite *ZskipListTestSuite) SetupTest() {
	suite.list = NewList()
}

func (suite *ZskipListTestSuite) TestNewList() {
	list := suite.list
	header := list.Header

	assert.Equal(suite.T(), 1, list.Level)
	assert.EqualValues(suite.T(), 0, header.Score)
	assert.Equal(suite.T(), consts.ZskipListMaxLevel, len(header.Level))
}

func (suite *ZskipListTestSuite) TestInsert() {
	list := suite.list
	obj := gredis.NewRawStringObject("abc")
	node, err := list.Insert(1, obj)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), obj, node.Obj)
	assert.EqualValues(suite.T(), 1, node.Score)
	assert.EqualValues(suite.T(), 1, list.Length)
}

func (suite *ZskipListTestSuite) TestGetRank() {
	list := suite.list

	obj := gredis.NewRawStringObject("abc")
	list.Insert(1, obj)

	obj2 := gredis.NewRawStringObject("xyb")
	list.Insert(2, obj2)

	obj3 := gredis.NewRawStringObject("def")
	list.Insert(1.3, obj3)

	assert.EqualValues(suite.T(), 1, list.GetRank(1, obj))
	assert.EqualValues(suite.T(), 3, list.GetRank(2, obj2))
	assert.EqualValues(suite.T(), 2, list.GetRank(1.3, obj3))
}

func (suite *ZskipListTestSuite) TestDelete() {
	list := suite.list

	obj := gredis.NewRawStringObject("abc")
	list.Insert(1, obj)

	obj2 := gredis.NewRawStringObject("xyb")
	list.Insert(2, obj2)

	assert.EqualValues(suite.T(), 2, list.Length)
	assert.True(suite.T(), list.Delete(2, obj2))
	assert.False(suite.T(), list.Delete(3, obj2))
	assert.EqualValues(suite.T(), 1, list.Length)
}

func (suite *ZskipListTestSuite) TestGetElementByRank() {
	list := suite.list

	obj := gredis.NewRawStringObject("abc")
	node1, _ := list.Insert(1, obj)

	obj2 := gredis.NewRawStringObject("xyb")
	node2, _ := list.Insert(2, obj2)

	find1 := list.GetElementByRank(1)
	find2 := list.GetElementByRank(2)

	assert.EqualValues(suite.T(), find2, node2)
	assert.EqualValues(suite.T(), find1, node1)
}

func (suite *ZskipListTestSuite) TestIsInRange() {
	list := suite.list

	obj := gredis.NewRawStringObject("abc")
	_, err := list.Insert(1, obj)
	suite.NoError(err)

	obj2 := gredis.NewRawStringObject("xyb")
	_, err = list.Insert(2, obj2)
	suite.NoError(err)

	assert.True(suite.T(), list.IsInRange(&ZRangeSpec{Min:0.3, Max:1}))
	assert.False(suite.T(), list.IsInRange(&ZRangeSpec{Min:22, Max:1}))
	assert.True(suite.T(), list.IsInRange(&ZRangeSpec{Min:1, Max:12}))
	assert.False(suite.T(), list.IsInRange(&ZRangeSpec{Min:2, Max:5, Minex:true}))
	assert.True(suite.T(), list.IsInRange(&ZRangeSpec{Min:2, Max:5}))
	assert.True(suite.T(), list.IsInRange(&ZRangeSpec{Min:1, Max:2}))
	assert.True(suite.T(), list.IsInRange(&ZRangeSpec{Min:1, Max:2}))
}

func TestZskipListTestSuite(t *testing.T) {
	suite.Run(t, new(ZskipListTestSuite))
}
