package dto

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
