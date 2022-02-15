package es

//一些需要用到的结构

type Json map[string]interface{}

type Sort struct {
	Field     string
	Ascending bool
}

type Page struct {
	PageNum  int
	PageSize int
}

type EsConnectID struct {
	EsConnectID int `json:"es_connect"`
}

type EsMapGetProperties struct {
	EsConnectID int    `json:"es_connect"`
	IndexName   string `json:"index_name"`
}

type UpdateMapping struct {
	EsConnect  int    `json:"es_connect"`
	IndexName  string `json:"index_name"`
	TypeName   string `json:"type_name"`
	Properties Json   `json:"properties"`
}

type TaskList struct {
	EsConnect int `json:"es_connect"`
}

type CancelTask struct {
	EsConnect int    `json:"es_connect"`
	TaskID    string `json:"task_id"`
}

type EsSnapshotInfo struct {
	EsConnect        int      `json:"es_connect"`
	SnapshotInfoList []string `json:"snapshot_info_list"`
}

type SnapshotCreateRepository struct {
	EsConnect              int    `json:"es_connect"`
	Repository             string `json:"name"`
	Type                   string `json:"type"`
	Location               string `json:"location"`
	Compress               string `json:"compress"`
	MaxRestoreBytesPerSec  string `json:"max_restore_bytes_per_sec"`
	MaxSnapshotBytesPerSec string `json:"max_snapshot_bytes_per_sec"`
	ChunkSize              string `json:"chunk_size"`
	Readonly               string `json:"readonly"`
}

type CleanupeRepository struct {
	EsConnect  int    `json:"es_connect"`
	Repository string `json:"name"`
}

type SnapshotDeleteRepository struct {
	EsConnect  int    `json:"es_connect"`
	Repository string `json:"name"`
}

type SnapshotList struct {
	EsConnect  int    `json:"es_connect"`
	Repository string `json:"repository"`
}

type CreateSnapshot struct {
	SnapshotName       string   `json:"snapshotName"`
	RepositoryName     string   `json:"repositoryName"`
	IndexList          []string `json:"indexList"`
	IgnoreUnavailable  *bool    `json:"ignore_unavailable"`
	IncludeGlobalState *bool    `json:"include_global_state"`
	Partial            *bool    `json:"partial"`
	Wait               *bool    `json:"wait"`
	EsConnect          int      `json:"es_connect"`
}

type DeleteSnapshot struct {
	SnapshotName   string `json:"snapshotName"`
	RepositoryName string `json:"repositoryName"`
	EsConnect      int    `json:"es_connect"`
}

type SnapshotDelete struct {
	EsConnect  int    `json:"es_connect"`
	Repository string `json:"repository"`
	Snapshot   string `json:"snapshot"`
}

type SnapshotDetail struct {
	EsConnect  int    `json:"es_connect"`
	Repository string `json:"repository"`
	Snapshot   string `json:"snapshot"`
}

type SnapshotRestore struct {
	SnapshotName       string   `json:"snapshotName"`
	RepositoryName     string   `json:"repositoryName"`
	IndexList          []string `json:"indexList"`
	IgnoreUnavailable  *bool    `json:"ignore_unavailable"`
	IncludeGlobalState *bool    `json:"include_global_state"`
	Partial            *bool    `json:"partial"`
	Wait               *bool    `json:"wait"`
	EsConnect          int      `json:"es_connect"`
	RenamePattern      string   `json:"rename_pattern"`
	RenameReplacement  string   `json:"rename_replacement"`
}

type SnapshotStatus struct {
	SnapshotName   string `json:"snapshot"`
	RepositoryName string `json:"repository"`
	EsConnect      int    `json:"es_connect"`
}

type EsConnect struct {
	Ip      string `json:"ip" db:"ip"`
	User    string `json:"user" db:"user"`
	Pwd     string `json:"pwd" db:"pwd"`
	Version int    `json:"version" db:"version"`
}

type EsCat struct {
	EsConnect        int    `json:"es_connect"`
	Cat              string `json:"cat"`
	IndexBytesFormat string `json:"index_bytes_format"`
}

type EsRest struct {
	EsConnect int    `json:"es_connect"`
	Method    string `json:"method"`
	Body      string `json:"body"`
	Path      string `json:"path"`
}

type EsOptimize struct {
	EsConnect int    `json:"es_connect"`
	IndexName string `json:"index_name"`
	Command   string `json:"command"`
}

type EsIndexInfo struct {
	EsConnect int    `json:"es_connect"`
	Settings  Json   `json:"settings"`
	IndexName string `json:"index_name"`
	Types     string `json:"types"`
}

type EsDocDeleteRowByID struct {
	EsConnect int    `json:"es_connect"`
	ID        string `json:"id"`
	IndexName string `json:"index_name"`
	Type      string `json:"type"`
}

type EsDocUpdateByID struct {
	EsConnect int    `json:"es_connect"`
	ID        string `json:"id"`
	JSON      Json   `json:"json"`
	Type      string `json:"type_name"`
	Index     string `json:"index"`
}

type EsMappingInfo struct {
	IndexNameList []string `json:"index_name_list"`
	EsConnect     int      `json:"es_connect"`
	Mappings      Json     `json:"mappings"`
	IndexName     string   `json:"index_name"`
}

type EsTaskInfo struct {
	EsConnect    int      `json:"es_connect"`
	TaskId       []string `json:"task_id"`
	Actions      []string `json:"actions"`
	NodeId       []string `json:"node_id"`
	ParentTaskId string   `json:"parent_task_id"`
}

type EsAliasInfo struct {
	EsConnect        int      `json:"es_connect"`
	Settings         Json     `json:"settings"`
	IndexName        string   `json:"index_name"`
	AliasName        string   `json:"alias_name"`
	NewAliasNameList []string `json:"new_alias_name_list"`
	NewIndexList     []string `json:"new_index_list"`
	Types            int      `json:"types"`
}

type EsReIndexInfo struct {
	EsConnect int `json:"es_connect"`
	UrlValues struct {
		Timeout             string `json:"timeout"`
		RequestsPerSecond   int    `json:"requests_per_second"`
		Slices              int    `json:"slices"`
		Scroll              string `json:"scroll"`
		WaitForActiveShards string `json:"wait_for_active_shards"`
		Refresh             string `json:"refresh"`
		WaitForCompletion   *bool  `json:"wait_for_completion"`
	} `json:"url_values"`
	Body map[string]interface{} `json:"body"`
}
