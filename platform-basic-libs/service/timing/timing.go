package timing

import (
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"log"
	"sync"
	"time"

	timing2 "github.com/lwl1989/timing"
)

type TaskFn func(data string) (err error)

var TimingFnMap = map[int]TaskFn{}

var TS *timing2.TaskScheduler
var once sync.Once

func GetTaskSchedulerInstance() *timing2.TaskScheduler {
	once.Do(func() {
		TS = timing2.GetTaskScheduler()
	})
	return TS
}

func AddTaskDb(taskTyp int, taskData string, taskID string, runtime int64) (err error) {
	gmTimedList := model.GmTimedList{}
	err = gmTimedList.AddTask(taskID, time.Unix(runtime, 0).Format(util.TimeFormat), taskData, taskTyp, runtime)
	if err != nil {
		return
	}
	return
}

func AddTask(taskTyp int, taskData string, taskID string, runtime int64) (err error) {
	gmTimedList := model.GmTimedList{}
	scheduler := GetTaskSchedulerInstance()

	task := &timing2.Task{
		Job: timing2.GetJob(func() {
			fn := TimingFnMap[taskTyp]
			log.Println("fn", fn)
			if err := fn(taskData); err != nil {
				panic(err)
			}
		}),
		Uuid:    taskID,
		RunTime: runtime,
	}

	task.GetJob().OnStart(func(reply timing2.Reply) {
		if err := gmTimedList.UpdateStatus(model.GmTimedListStart, "开始运行中...", taskID); err != nil {
			logs.Logger.Sugar().Errorf("err", err)
		}
	})

	task.GetJob().OnError(func(reply timing2.Reply) {
		if err := gmTimedList.UpdateStatus(model.GmTimedListErr, reply.Msg, taskID); err != nil {
			logs.Logger.Sugar().Errorf("err", err)
		}
	})

	task.GetJob().OnFinish(func(reply timing2.Reply) {
		if err := gmTimedList.UpdateStatus(model.GmTimedListSucc, reply.Msg, taskID); err != nil {
			logs.Logger.Sugar().Errorf("err", err)
		}
	})

	go scheduler.AddTask(task)

	return
}
