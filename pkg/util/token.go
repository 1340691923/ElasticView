package util

import (
	"github.com/google/uuid"
	"sync"
)

var TokenBucket sync.Map

func GetUUid() string {
	return uuid.New().String()
}
