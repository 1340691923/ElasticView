package util

import (
	"github.com/sony/sonyflake"
	"log"
	"strconv"
	"sync"
)

var TokenBucket sync.Map

func GetUUid() string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Println("err", err)
	}
	return strconv.Itoa(int(id))
}
