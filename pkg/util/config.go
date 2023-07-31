package util

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// LoadJSONConfig 读取配置文件 json格式
func LoadJSONConfig(filename string, v interface{}) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func JoinInt(s []int, sp string) string {
	var tmp = make([]string, 0, len(s))
	for i, _ := range s {
		tmp = append(tmp, strconv.Itoa(s[i]))
	}
	return strings.Join(tmp, sp)
}
