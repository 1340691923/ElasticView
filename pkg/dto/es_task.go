package dto

type TaskList struct {
	EsConnect int `json:"es_connect"`
}

type CancelTask struct {
	EsConnect int    `json:"es_connect"`
	TaskID    string `json:"task_id"`
}

type EsTaskInfo struct {
	EsConnect    int      `json:"es_connect"`
	TaskId       []string `json:"task_id"`
	Actions      []string `json:"actions"`
	NodeId       []string `json:"node_id"`
	ParentTaskId string   `json:"parent_task_id"`
}
