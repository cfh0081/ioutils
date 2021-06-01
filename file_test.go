package ioutils

import (
	"fmt"
	"io/ioutil"
	"os"

	//	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	// ioutil.TempDir 不指定第一个参数，则默认使用os.TempDir()目录
	tmepDir := os.TempDir()
	fmt.Println(tmepDir)

	tempDir, err := ioutil.TempDir("", "fileTest") //在DIR目录下创建tmp为目录名前缀的目录，DIR必须存在，否则创建不成功
	if err != nil {
		fmt.Println("临时目录创建失败")
		return
	}
	fmt.Println(tempDir) //生成的目录名为tmpXXXXX，XXXXX为一个随机数

	defer os.RemoveAll(tempDir)

	testFilePath := tempDir + string(os.PathSeparator) + "existsTest.txt"
	r := Exists(testFilePath)
	if r {
		t.Errorf("Exists(%s) failed. Got %t, expected false.", testFilePath, r)
		return
	}

	//	path, _ := os.Getwd()
	//	fmt.Println(path)
	//	fmt.Println(filepath.Abs(testFilePath))

	//	file, error := ioutil.TempFile("DIR", "tmp") //在DIR目录下创建tmp为文件名前缀的文件，获得file文件指针，DIR必须存在，否则创建不成功
	file, error := os.Create(testFilePath)
	if error != nil {
		fmt.Println("文件创建失败")
		return
	}
	file.WriteString("insert into file") //利用file指针的WriteString()写入内容
	file.Close()

	r = Exists(testFilePath)
	if !r {
		t.Errorf("Exists(%s) failed. Got %t, expected true.", testFilePath, r)
	}
}

func TestPureName(t *testing.T) {
	name := "what.bat"
	expect := "what"
	pure := PureName(name)
	r := pure == expect
	if !r {
		t.Errorf("PureName(%s) failed. Got %s, expected %s.", name, pure, expect)
	}

	name = "name"
	expect = "name"
	pure = PureName(name)
	r = pure == expect
	if !r {
		t.Errorf("PureName(%s) failed. Got %s, expected %s.", name, pure, expect)
	}
}
