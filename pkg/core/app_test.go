package core

import (
	"log"
	"testing"
)

func Test_register(t *testing.T) {
	Register(MaxLevel,"test1", func() (deferFn func(), err error) {
		log.Println("测试1")
		return func() {
			log.Println("结束测试1")
		},nil
	})
	Register(MiddleLevel,"test2", func() (deferFn func(), err error) {
		log.Println("测试2")
		return func() {
			log.Println("结束测试2")
		},nil
	})
	Register(MinLevel,"test3", func() (deferFn func(), err error) {
		log.Println("测试3")
		return func() {
			log.Println("结束测试3")
		},nil
	})
	Run()
	defer Stop()
}
