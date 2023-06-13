package response

import (
	"github.com/1340691923/ElasticView/model"
	"strings"
)

type SearchConfig struct {
	ID         int      `db:"id" json:"id"`
	IndexName  string   `db:"index_name" json:"indexName"`
	Remark     string   `db:"remark" json:"remark"`
	InputCols  []string `json:"input_cols" db:"input_cols"`
	OutputCols []string `json:"output_cols" db:"output_cols"`

	EsConnect int    `db:"es_connect" json:"-"`
	Updated   string `db:"updated" json:"updated"`
	Created   string `db:"created" json:"created"`
	Limit     int    `json:"-" db:"-"`
	Page      int    `json:"-" db:"-"`
}

func ToSearchConfig(list []model.SearchConfig) []SearchConfig {
	res := []SearchConfig{}
	for _, v := range list {
		res = append(res, SearchConfig{
			ID:         v.ID,
			IndexName:  v.IndexName,
			Remark:     v.Remark,
			InputCols:  strings.Split(v.InputCols, ","),
			OutputCols: strings.Split(v.OutputCols, ","),
			EsConnect:  v.EsConnect,
			Updated:    v.Updated,
			Created:    v.Created,
			Limit:      v.Limit,
			Page:       v.Page,
		})
	}
	return res
}
