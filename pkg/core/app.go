package core

import (
	"fmt"
	"log"
	"sort"
)

type LevelAndRegiter struct {
	level Level
	tag   string
	fn    InitFnObserver
}

func NewLevelAndRegiter(level Level, tag string, fn InitFnObserver) *LevelAndRegiter {
	return &LevelAndRegiter{level: level, tag: tag, fn: fn}
}

var (
	registerList = []*LevelAndRegiter{}
	deferFnList  = []func(){}
)

type Level int

const (
	MaxLevel Level = iota
	MiddleLevel
	MinLevel
	LastLevel
)

type InitFnObserver func() (deferFn func(), err error)

func Register(level Level, tag string, fn InitFnObserver) {
	registerList = append(registerList, NewLevelAndRegiter(level, tag, fn))
}

func Run() {
	sort.Slice(registerList, func(i, j int) bool {
		return registerList[i].level < registerList[j].level
	})
	for _, register := range registerList {
		deferFn, err := register.fn()
		if err != nil {
			panic(fmt.Sprintf("%s组件初始化失败:%s", register.tag, err.Error()))
		} else {
			log.Println(fmt.Sprintf("%s组件初始化成功", register.tag))
		}
		deferFnList = append(deferFnList, deferFn)
	}
}

func Stop() {
	reverseArray(deferFnList)
	for _, fn := range deferFnList {
		fn()
	}
}

func reverseArray(arr []func()) []func() {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
