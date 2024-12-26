package util

import (
	"fmt"
	"github.com/pkg/errors"
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
	var err error
	// 遍历原文件夹内部所有item
	items, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, item := range items {

		// 文件
		if !item.IsDir() {
			err = cpoyFile2(path.Join(src, item.Name()), path.Join(dest, item.Name()))
			if err != nil {
				return errors.WithStack(err)
			}
			continue
		}

		// 目录
		err = os.Mkdir(path.Join(dest, item.Name()), os.ModePerm)
		if err != nil {
			return errors.WithStack(err)
		}
		// 递归
		err = DirCopy(path.Join(src, item.Name()), path.Join(dest, item.Name()))
		if err != nil {
			return errors.WithStack(err)
		}
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

// MoveDir moves a folder and its contents to a new location.
func MoveDir(srcDir, destDir string) error {
	// Ensure the source directory exists
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("source directory error: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// Create the destination directory if it doesn't exist
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	}

	// Walk through the source directory and move files
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create the destination directory
			return os.MkdirAll(destPath, os.ModePerm)
		}

		// Move the file
		return moveFile(path, destPath)
	})
	if err != nil {
		return err
	}

	// Remove the source directory
	return os.RemoveAll(srcDir)
}

// moveFile moves a single file to a new location
func moveFile(srcFile, destFile string) error {
	// Open the source file
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the destination file
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	// Copy the file contents
	if _, err := io.Copy(dest, src); err != nil {
		return err
	}

	// Remove the source file
	return os.Remove(srcFile)
}
