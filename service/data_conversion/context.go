package data_conversion

import (
	"context"
	"errors"
	"sync"
)

var TS *Task
var once sync.Once

type Task struct {
	cancelMap map[int]context.CancelFunc
	lock      *sync.RWMutex
}

func (this *Task) SetCancelFunc(id int, cancel context.CancelFunc) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.cancelMap[id] = cancel
}

func (this *Task) CancelById(id int) (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.cancelMap[id]
	if !ok {
		return errors.New("没有找到该任务id")
	}
	this.cancelMap[id]()
	return nil
}

func newTask() *Task {
	return &Task{
		cancelMap: map[int]context.CancelFunc{},
		lock:      new(sync.RWMutex),
	}
}

func GetTaskInstance() *Task {
	once.Do(func() {
		TS = newTask()
	})
	return TS
}
