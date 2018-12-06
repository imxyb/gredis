package dict

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"github.com/stretchr/testify/assert"
)

type DictTestSuite struct {
	suite.Suite
	Dict *Dict
}

func (suite *DictTestSuite) SetupTest() {
	suite.Dict = NewDict()
}

func (suite *DictTestSuite) TestAdd() {
	suite.Dict.Add(1, 2)
	assert.EqualValues(suite.T(), 2, suite.Dict.FetchKey(1))
}

func (suite *DictTestSuite) TestReplace() {
	suite.Dict.Add(1, "a")

	assert.EqualValues(suite.T(), "a", suite.Dict.FetchKey(1))
	suite.Dict.Replace(1, "aa")
	assert.EqualValues(suite.T(), "aa", suite.Dict.FetchKey(1))

	suite.Dict.Replace(2, "b")
	assert.EqualValues(suite.T(), "b", suite.Dict.FetchKey(2))
}

func (suite *DictTestSuite) TestFetchKey() {
	suite.Dict.Add(1, "a")
	assert.EqualValues(suite.T(), "a", suite.Dict.FetchKey(1))
	assert.Nil(suite.T(), suite.Dict.FetchKey(11))
}

func (suite *DictTestSuite) TestGetRandomKey() {
	suite.Dict.Add(1, "a")
	suite.Dict.Add(2, "b")

	result := suite.Dict.GetRandomKey()
	assert.Condition(suite.T(), func() (success bool) {
		for key, value := range result {
			k := key.(int)
			v := value.(string)

			if (k == 1 || k == 2) && (v == "a" || v == "b") {
				success = true
				break
			}

			success = false
		}
		return
	})
}

func (suite *DictTestSuite) TestDel() {
	suite.Dict.Add(1, "a")
	suite.Dict.Add(2, "b")

	suite.Dict.Delete(1)
	assert.Nil(suite.T(), suite.Dict.FetchKey(1))
	assert.EqualValues(suite.T(), "b", suite.Dict.FetchKey(2))
}

func TestDictTestSuite(t *testing.T) {
	suite.Run(t, new(DictTestSuite))
}
