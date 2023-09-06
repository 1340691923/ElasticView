package util

import (
	"github.com/google/uuid"
	"github.com/sony/sonyflake"
	"log"
	"strconv"
	"sync"
	"time"
)

var TokenBucket sync.Map

func GetUUid() string {
	return uuid.New().String()
}

var UUIDBucket sync.Map

func GetUUid2() string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Println("err", err)
	}
	uuid := strconv.Itoa(int(id))
	if _, ok := UUIDBucket.Load(uuid); ok {
		time.Sleep(100 * time.Millisecond)
		return GetUUid2()
	}
	UUIDBucket.Store(uuid, struct{}{})
	return uuid
}
