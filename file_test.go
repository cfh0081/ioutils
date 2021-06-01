package ioutils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	//	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	// ioutil.TempDir 不指定第一个参数，则默认使用os.TempDir()目录
	// tmepDir := os.TempDir()
	// fmt.Println(tmepDir)

	tempDir, err := ioutil.TempDir("", "fileTest") //在DIR目录下创建tmp为目录名前缀的目录，DIR必须存在，否则创建不成功
	assert.Nil(t, err)
	// fmt.Println(tempDir) //生成的目录名为tmpXXXXX，XXXXX为一个随机数
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "existsTest.txt")
	r := Exists(testFilePath)
	assert.False(t, r)

	//	path, _ := os.Getwd()
	//	fmt.Println(path)
	//	fmt.Println(filepath.Abs(testFilePath))

	//	file, error := ioutil.TempFile("DIR", "tmp") //在DIR目录下创建tmp为文件名前缀的文件，获得file文件指针，DIR必须存在，否则创建不成功
	file, err := os.Create(testFilePath)
	assert.Nil(t, err)
	_, err = file.WriteString("insert into file") //利用file指针的WriteString()写入内容
	assert.Nil(t, err)
	file.Close()

	r = Exists(testFilePath)
	assert.True(t, r)
}

func TestPureName(t *testing.T) {
	name := "what.bat"
	expected := "what"
	pure := PureName(name)
	assert.Equal(t, expected, pure)

	name = "name"
	expected = "name"
	pure = PureName(name)
	assert.Equal(t, expected, pure)
}

func TestIsDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "dirTest") //在DIR目录下创建tmp为目录名前缀的目录，DIR必须存在，否则创建不成功
	assert.Nil(t, err)

	defer os.RemoveAll(tempDir)

	testFileDir := filepath.Join(tempDir, "sub")
	r := IsDir(testFileDir)
	assert.False(t, r)

	file, err := os.Create(testFileDir)
	assert.Nil(t, err)
	_, err = file.WriteString("insert into file") //利用file指针的WriteString()写入内容
	assert.Nil(t, err)
	file.Close()
	r = IsDir(testFileDir)
	assert.False(t, r)
	os.Remove(testFileDir)

	err = os.MkdirAll(testFileDir, os.ModePerm)
	assert.Nil(t, err)

	r = IsDir(testFileDir)
	assert.True(t, r)
}
