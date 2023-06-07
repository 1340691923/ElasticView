package util

import "encoding/json"

func IsJson(str string) bool {
	maps := make(map[string]interface{})

	err := json.Unmarshal(Str2bytes(str), &maps)
	if err != nil {
		return false
	}
	return true
}
