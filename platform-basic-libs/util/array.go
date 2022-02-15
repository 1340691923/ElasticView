// 辅助工具层
package util

import (
	"errors"
	"strconv"
	"strings"
	"unsafe"
)

func InArr(array []int, column int) bool {
	i := 0
	for i < len(array) {
		if array[i] == column {
			return true
		}
		i++
	}
	return false
}

func InstrArr(array []string, column string) bool {
	i := 0
	for i < len(array) {
		if array[i] == column {
			return true
		}
		i++
	}
	return false
}

func InMap(maps map[string]int, column string) (ok bool) {
	_, ok = maps[column]
	return
}

//替换string与byte转换时性能损耗的代码  ！！！ 只可用于不可修改字符串变量
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// SplitInt 分割字符串并转为INT
func SplitInt(s, sep string) ([]int, error) {
	var tmpStr = strings.Split(s, sep)
	if len(tmpStr) == 0 {
		return nil, errors.New("no split int")
	}
	var err error
	var tmp int
	var rlt []int
	for i := range tmpStr {
		tmp, err = strconv.Atoi(tmpStr[i])
		if err != nil {
			return nil, err
		}
		rlt = append(rlt, tmp)
	}
	return rlt, nil
}
