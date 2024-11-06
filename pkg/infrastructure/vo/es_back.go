package vo

type SnapshotRepositoryList struct {
	List     []Snashot                     `json:"list"`
	Res      map[string]SnapshotRepository `json:"res"`
	PathRepo []interface{}                 `json:"pathRepo"`
}
