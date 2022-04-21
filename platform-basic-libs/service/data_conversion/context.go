package data_conversion

import (
	"context"
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
