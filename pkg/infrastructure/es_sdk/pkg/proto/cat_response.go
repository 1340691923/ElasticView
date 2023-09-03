package proto

type CatIndex struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	Uuid         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}
