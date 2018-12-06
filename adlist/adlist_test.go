package adlist

import (
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AdListTestSuite struct {
	suite.Suite
	List *List
}

func (suite *AdListTestSuite) SetupTest() {
	suite.List = NewList()
}

func (suite *AdListTestSuite) TestAddNodeHead() {
	list := suite.List

	list.AddNodeHead(123)
	assert.Equal(suite.T(), 123, list.Head.Value)
	assert.Equal(suite.T(), 123, list.Tail.Value)
	assert.EqualValues(suite.T(), 1, list.Len)

	list.AddNodeHead(456)
	assert.Equal(suite.T(), 456, list.Head.Value)
	assert.Equal(suite.T(), 123, list.Tail.Value)
	assert.EqualValues(suite.T(), 2, list.Len)

	list.AddNodeHead(789)
	assert.Equal(suite.T(), 789, list.Head.Value)
	assert.Equal(suite.T(), 123, list.Tail.Value)
	assert.EqualValues(suite.T(), 3, list.Len)
}

func (suite *AdListTestSuite) TestAddNodeTail() {
	list := suite.List

	list.AddNodeTail(123)
	assert.Equal(suite.T(), 123, list.Head.Value)
	assert.Equal(suite.T(), 123, list.Tail.Value)
	assert.EqualValues(suite.T(), 1, list.Len)

	list.AddNodeTail(456)
	assert.Equal(suite.T(), 123, list.Head.Value)
	assert.Equal(suite.T(), 456, list.Tail.Value)
	assert.EqualValues(suite.T(), 2, list.Len)

	list.AddNodeTail(789)
	assert.Equal(suite.T(), 123, list.Head.Value)
	assert.Equal(suite.T(), 789, list.Tail.Value)
	assert.EqualValues(suite.T(), 3, list.Len)
}

func (suite *AdListTestSuite) TestSearchKey() {
	list := suite.List

	list.AddNodeHead(123)
	list.AddNodeHead(456)

	assert.Equal(suite.T(), 123, list.SearchKey(123).Value)
	assert.Nil(suite.T(), list.SearchKey(1234))
	assert.Equal(suite.T(), 456, list.SearchKey(456).Value)
}

func (suite *AdListTestSuite) TestInsertNode() {
	list := suite.List

	list.AddNodeHead(123)
	list.AddNodeHead(456)

	node1 := list.SearchKey(123)
	node2 := list.SearchKey(456)

	list.InsertNode(node2, 789, 0)
	assert.Equal(suite.T(), 789, list.Head.Value)
	assert.EqualValues(suite.T(), 3, list.Len)

	list.InsertNode(node1, 888, 1)
	assert.Equal(suite.T(), 888, list.Tail.Value)
	assert.EqualValues(suite.T(), 4, list.Len)

	iter1 := NewListIter(list, AL_START_HEAD)

	var result1 []interface{}
	for {
		node := iter1.Next()
		if node == nil {
			break
		}

		result1 = append(result1, node.Value)
	}

	e1 := []interface{}{789, 456, 123, 888}
	assert.Equal(suite.T(), e1, result1)

	iter2 := NewListIter(list, AL_START_TAIL)

	var result2 []interface{}
	for {
		node := iter2.Next()
		if node == nil {
			break
		}

		result2 = append(result2, node.Value)
	}

	e2 := []interface{}{888, 123, 456, 789}
	assert.Equal(suite.T(), e2, result2)
}

func (suite *AdListTestSuite) TestIndex() {
	list := suite.List

	list.AddNodeTail(123)
	list.AddNodeTail(456)
	list.AddNodeTail(789)

	assert.Equal(suite.T(), 123, list.Index(0).Value)
	assert.Equal(suite.T(), 456, list.Index(1).Value)
	assert.Equal(suite.T(), 789, list.Index(2).Value)

	assert.Equal(suite.T(), 789, list.Index(-1).Value)
	assert.Equal(suite.T(), 456, list.Index(-2).Value)
	assert.Equal(suite.T(), 123, list.Index(-3).Value)
}

func (suite *AdListTestSuite) TestDelNode() {
	list := suite.List

	list.AddNodeTail(123)
	list.AddNodeTail(456)
	list.AddNodeTail(789)

	node1 := list.SearchKey(123)
	node2 := list.SearchKey(456)
	node3 := list.SearchKey(789)

	list.DelNode(node1)
	assert.Equal(suite.T(), 456, list.Head.Value)
	assert.Equal(suite.T(), 789, list.Tail.Value)
	assert.EqualValues(suite.T(), 2, list.Len)

	list.DelNode(node2)
	assert.Equal(suite.T(), 789, list.Head.Value)
	assert.EqualValues(suite.T(), 1, list.Len)

	list.DelNode(node3)
	assert.Nil(suite.T(), list.Head)
	assert.EqualValues(suite.T(), 0, list.Len)
}

func (suite *AdListTestSuite) TestListDup() {
	origin := suite.List
	origin.AddNodeTail(123)
	origin.AddNodeTail(456)

	cpy := ListDup(origin)
	assert.Equal(suite.T(), 123, cpy.SearchKey(123).Value)
	assert.Equal(suite.T(), 456, cpy.SearchKey(456).Value)
}

func TestAdListTestSuite(t *testing.T) {
	suite.Run(t, new(AdListTestSuite))
}


