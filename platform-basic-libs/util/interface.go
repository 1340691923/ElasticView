package util

import "reflect"

// interface{}转为 []interface{}
func CreateAnyTypeSlice(slice interface{}) ([]interface{}, bool) {
	val, ok := isSlice(slice)

	if !ok {
		return nil, false
	}

	sliceLen := val.Len()

	out := make([]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

// 判断是否为slcie数据
func isSlice(arg interface{}) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return
}

