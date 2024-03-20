package file_demo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// 文件的基本操作

func TestFile(t *testing.T) {

	// 读写文件
	filepath, _ := os.Getwd()
	fmt.Println(filepath)
	file, err := os.Open("testdata/my_file.txt") // 在进行写操作时，会出现 bad file descriptor -> 不可写 打开模式 权限
	//file, err := os.OpenFile("testdata/my_file.txt", os.O_RDWR|os.O_CREATE, 0666)
	require.NoError(t, err)
	data := make([]byte, 128)
	n, err := file.Read(data)
	fmt.Println(n)
	require.NoError(t, err)

	//file.WriteString("\n")
	//n, err = file.WriteString("hello my_file")
	//fmt.Println(n)
	//fmt.Println(err)
	//require.NoError(t, err)

	file.Close()

	file, err = os.OpenFile("testdata/my_file.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	require.NoError(t, err)
	file.WriteString("\n")
	n, err = file.WriteString("hello my_file")
	fmt.Println(n)
	fmt.Println(err)
	require.NoError(t, err)

	file.Close()

	// 新建文件

	file, err = os.Create("testdata/my_file_copy.txt")
	require.NoError(t, err)
	n, err = file.WriteString("hello tomato")
	fmt.Println(n)
	require.NoError(t, err)
}
