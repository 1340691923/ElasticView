package dto

type EsAliasInfo struct {
	EsConnect        int      `json:"es_connect"`
	Settings         Json     `json:"settings"`
	IndexName        string   `json:"index_name"`
	AliasName        string   `json:"alias_name"`
	NewAliasNameList []string `json:"new_alias_name_list"`
	NewIndexList     []string `json:"new_index_list"`
	Types            int      `json:"types"`
}
