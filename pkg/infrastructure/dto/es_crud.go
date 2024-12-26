package dto

type CrudFilter struct {
	Relation  AnalysisFilter `json:"relation"`
	SortList  []SortStruct   `json:"sort_list"`
	EsConnect int            `json:"es_connect"`
	IndexName string         `json:"index_name"`
	Page      int            `json:"page"`
	Limit     int            `json:"limit"`
}

type AnalysisFilter struct {
	FilterType string `json:"filterType"`
	Filts      []struct {
		FilterType string `json:"filterType"`
		Filts      []struct {
			ColumnName string      `json:"columnName"`
			Comparator string      `json:"comparator"`
			FilterType string      `json:"filterType"`
			Ftv        interface{} `json:"ftv"`
		} `json:"filts,omitempty"`
		Relation   string      `json:"relation,omitempty"`
		ColumnName string      `json:"columnName,omitempty"`
		Comparator string      `json:"comparator,omitempty"`
		Ftv        interface{} `json:"ftv,omitempty"`
	} `json:"filts"`
	Relation string `json:"relation"`
}

type SortStruct struct {
	Col      string `json:"col"`
	SortRule string `json:"sortRule"`
}
