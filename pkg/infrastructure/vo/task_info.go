package vo

type TaskInfo struct {
	Node               string            `json:"node"`
	Id                 int64             `json:"id"` // the task id (yes, this is a long in the Java source)
	Type               string            `json:"type"`
	Action             string            `json:"action"`
	Status             interface{}       `json:"status"`      // has separate implementations of Task.Status in Java for reindexing, replication, and "RawTaskStatus"
	Description        interface{}       `json:"description"` // same as Status
	StartTime          string            `json:"start_time"`
	StartTimeInMillis  int64             `json:"start_time_in_millis"`
	RunningTime        string            `json:"running_time"`
	RunningTimeInNanos int64             `json:"running_time_in_nanos"`
	Cancellable        bool              `json:"cancellable"`
	ParentTaskId       string            `json:"parent_task_id"` // like "YxJnVYjwSBm_AUbzddTajQ:12356"
	Headers            map[string]string `json:"headers"`
}

type Tasks struct {
	Nodes map[string]TaskNode `json:"nodes"`
}

type TaskNode struct {
	Name             string   `json:"name"`
	TransportAddress string   `json:"transport_address"`
	Host             string   `json:"host"`
	Ip               string   `json:"ip"`
	Roles            []string `json:"roles"`
	Attributes       struct {
		MlMachineMemory string `json:"ml.machine_memory"`
		XpackInstalled  string `json:"xpack.installed"`
		MlMaxOpenJobs   string `json:"ml.max_open_jobs"`
		MlEnabled       string `json:"ml.enabled"`
	} `json:"attributes"`
	Tasks map[string]TaskInfo `json:"tasks"`
}
