package util

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func WriteVersionGoFile(content string) (err error) {
	err = ioutil.WriteFile(`pkg\infrastructure\config\version.go`, []byte(content), 0777)
	return
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func DirCopy(src string, dest string) error {
	// 遍历原文件夹内部所有item
	items, _ := ioutil.ReadDir(src)
	for _, item := range items {

		// 文件
		if !item.IsDir() {
			cpoyFile2(path.Join(src, item.Name()), path.Join(dest, item.Name()))
			continue
		}

		// 目录
		os.Mkdir(path.Join(dest, item.Name()), os.ModePerm)
		// 递归
		DirCopy(path.Join(src, item.Name()), path.Join(dest, item.Name()))
	}

	return nil
}

func cpoyFile2(src, dest string) error {
	// open src readonly
	srcFp, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFp.Close()

	// create dest
	dstFp, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dstFp.Close()

	// copy
	_, err = io.Copy(dstFp, srcFp)
	return err
}
