package proto

type AliasAddAction struct {
	Add AliasAdd `json:"add"`
}

type AliasAdd struct {
	Indices []string `json:"indices"`
	Alias   string   `json:"alias"`
}

type AliasAction struct {
	Actions []AliasAddAction `json:"actions"`
}
