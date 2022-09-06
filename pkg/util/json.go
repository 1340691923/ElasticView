package util

import (
	jsoniter "github.com/json-iterator/go"
)

func IsJson(str string) bool {
	maps := make(map[string]interface{})
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(Str2bytes(str), &maps)
	if err != nil {
		return false
	}
	return true
}
