package sds

import (
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SdsTestSuite struct {
	suite.Suite
	Sds *Sds
}

func (suite *SdsTestSuite) SetupTest() {
	suite.Sds = NewSds("hello")
}

func (suite *SdsTestSuite) TestSLen() {
	assert.Equal(suite.T(), 5, suite.Sds.SLen())
}

func (suite *SdsTestSuite) TestAvail() {
	assert.Equal(suite.T(), 0, suite.Sds.Avail())
}

func (suite *SdsTestSuite) TestMakeRoomFor() {
	suite.Sds.MakeRoomFor(10)
	assert.Equal(suite.T(), (suite.Sds.SLen() + 10)*2, cap(suite.Sds.buf))

	suite.Sds.MakeRoomFor(1024 * 1024)
	assert.Equal(suite.T(), (suite.Sds.SLen() + 1024 * 1024) + 1024 * 1024, cap(suite.Sds.buf))
}

func (suite *SdsTestSuite) TestDupSds() {
	newSds := DupSds(suite.Sds)
	assert.Equal(suite.T(), "hello", string(newSds.buf))
	assert.Equal(suite.T(), 5, newSds.SLen())
}

func (suite *SdsTestSuite) TestClear() {
	suite.Sds.Clear()
	assert.Equal(suite.T(), "", string(suite.Sds.buf))
	assert.Equal(suite.T(), 0, suite.Sds.SLen())
	assert.Equal(suite.T(), 5, suite.Sds.Avail())
}

func (suite *SdsTestSuite) TestSdsCat() {
	suite.Sds.StrCat("world")
	assert.Equal(suite.T(), "helloworld", string(suite.Sds.buf))
	assert.Equal(suite.T(), 10, suite.Sds.SLen())
	assert.Equal(suite.T(), 10, suite.Sds.Avail())
}

func (suite *SdsTestSuite) TestSdsCatSds() {
	newSds := NewSds("world")
	suite.Sds.StrCatSds(*newSds)
	assert.Equal(suite.T(), "helloworld", string(suite.Sds.buf))
	assert.Equal(suite.T(), 10, suite.Sds.SLen())
	assert.Equal(suite.T(), 10, suite.Sds.Avail())
}

func (suite *SdsTestSuite) TestCopy() {
	suite.Sds.Copy("imxyb")
	assert.Equal(suite.T(), "imxyb", string(suite.Sds.buf))
	assert.Equal(suite.T(), 5, suite.Sds.SLen())
	assert.Equal(suite.T(), 0, suite.Sds.Avail())

	suite.Sds.Copy("helloworld")
	assert.Equal(suite.T(), "helloworld", string(suite.Sds.buf))
	assert.Equal(suite.T(), 10, suite.Sds.SLen())
	assert.Equal(suite.T(), 10, suite.Sds.Avail())
}

func (suite *SdsTestSuite) TestRange() {
	err := suite.Sds.Range(1, -1)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "ello", string(suite.Sds.buf))
	assert.Equal(suite.T(), 4, suite.Sds.SLen())
	assert.Equal(suite.T(), 1, suite.Sds.Avail())

	err = suite.Sds.Range(10, 10)
	assert.Error(suite.T(), err)
}

func (suite *SdsTestSuite) TestTrim() {
	suite.Sds.Trim("he")
	assert.Equal(suite.T(), "llo", string(suite.Sds.buf))
	assert.Equal(suite.T(), 3, suite.Sds.SLen())
	assert.Equal(suite.T(), 2, suite.Sds.Avail())
}

func (suite *SdsTestSuite) TestCmpSds() {
	current := suite.Sds

	s1 := NewSds("helld")
	s2 := NewSds("hello")
	s3:= NewSds("hellz")

	assert.Equal(suite.T(), 1, CmpSds(current, s1))
	assert.Equal(suite.T(), 0, CmpSds(current, s2))
	assert.Equal(suite.T(), -1, CmpSds(current, s3))
}

func TestSdsTestSuite(t *testing.T) {
	suite.Run(t, new(SdsTestSuite))
}
