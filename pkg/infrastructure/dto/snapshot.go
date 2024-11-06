package dto

type EsSnapshotInfo struct {
	EsConnect        int      `json:"es_connect"`         //es连接id
	SnapshotInfoList []string `json:"snapshot_info_list"` //存储库
}

type SnapshotCreateRepository struct {
	EsConnect              int    `json:"es_connect"`                 //es连接id
	Repository             string `json:"name"`                       //存储库名
	Type                   string `json:"type"`                       //类型 fs/url
	Location               string `json:"location"`                   //存储位置
	Compress               string `json:"compress"`                   //是否压缩 true/false
	MaxRestoreBytesPerSec  string `json:"max_restore_bytes_per_sec"`  //节点恢复速率
	MaxSnapshotBytesPerSec string `json:"max_snapshot_bytes_per_sec"` //每个节点快照速率
	ChunkSize              string `json:"chunk_size"`                 //大文件分解块大小
	Readonly               string `json:"readonly"`                   //是否只读
}

type CleanupeRepository struct {
	EsConnect  int    `json:"es_connect"` //es连接id
	Repository string `json:"name"`       //存储库名
}

type SnapshotDeleteRepository struct {
	EsConnect  int    `json:"es_connect"` //es连接id
	Repository string `json:"name"`       //存储库名
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
