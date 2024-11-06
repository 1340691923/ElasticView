// 辅助工具层
package util

import (
	"github.com/pkg/errors"
	"reflect"
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
	l := len(array)
	for i < l {
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

func RemoveRepeatedElement(arr []string) (newArr []interface{}) {
	newArr = make([]interface{}, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func Bytes2str(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

// s2b converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func Str2bytes(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
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
