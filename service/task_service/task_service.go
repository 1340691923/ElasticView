package task_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/1340691923/ElasticView/es_sdk/pkg"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/vo"
	"github.com/tidwall/gjson"
)

type TaskService struct {
	esClient pkg.EsI
}

func NewTaskService(esClient pkg.EsI) *TaskService {
	return &TaskService{esClient: esClient}
}

func (this TaskService) TaskList(ctx context.Context) (res map[string]vo.TaskInfo, err error) {

	resp, err := this.esClient.TaskList(ctx)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	var tasksListResponse vo.Tasks
	err = json.Unmarshal(resp.ResByte(), &tasksListResponse)
	if err != nil {
		return
	}

	res = map[string]vo.TaskInfo{}

	for _, node := range tasksListResponse.Nodes {
		for taskId, taskInfo := range node.Tasks {
			res[taskId] = taskInfo
		}
	}

	return
}

func (this TaskService) Cancel(ctx context.Context, cancelTask *escache.CancelTask) (err error) {
	res, err := this.esClient.TasksCancel(ctx, cancelTask.TaskID)
	if err != nil {
		return
	}

	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	if gjson.GetBytes(res.ResByte(), "node_failures").Exists() {
		err = errors.New(string(res.ResByte()))
		return
	}

	return
}
