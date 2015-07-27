package util

import (
	"os"
	"path"
	"testing"

	. "gopkg.in/check.v1"
)

const (
	testPath = "./test_path"
)

func Test(t *testing.T) {
	TestingT(t)
}

type TestUtilSuite struct{}

var _ = Suite(&TestUtilSuite{})

func (s *TestUtilSuite) SetUpSuite(c *C) {
	err := os.MkdirAll(testPath, 0755)
	c.Assert(err, IsNil)
}

func (s *TestUtilSuite) TearDownSuite(c *C) {
	err := os.RemoveAll(testPath)
	c.Assert(err, IsNil)
}

func (s *TestUtilSuite) TestReadFile(c *C) {
	fileName := path.Join(testPath, "test_r.txt")
	buf := []byte("123xx..end")

	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	c.Assert(err, IsNil)

	_, err = fp.WriteAt(buf, 0)
	c.Assert(err, IsNil)

	ret, err := ReadFile(fileName)
	c.Assert(err, IsNil)
	c.Assert(string(ret), Equals, string(ret))
}

func (s *TestUtilSuite) TestCreateLog(c *C) {
	_, err := CreateLog(testPath+"/log", "create_successful")
	c.Assert(err, IsNil)
}

func (s *TestUtilSuite) TestCheckIsEmpty(c *C) {
	c.Assert(CheckIsEmpty("a", "..."), Equals, false)

	c.Assert(CheckIsEmpty("", ""), Equals, true)
}

func (s *TestUtilSuite) TestContains(c *C) {
	str := ""
	strs := []string{}
	c.Assert(Contains(str, strs), Equals, false)

	strs = append(strs, "123", "xyz")
	c.Assert(Contains(str, strs), Equals, false)

	strs = append(strs, "", "")
	c.Assert(Contains(str, strs), Equals, true)
}
