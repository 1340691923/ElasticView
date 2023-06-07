package request

import (
	"encoding/json"
	"errors"
)

// GmRoleModel
type GmRoleModel struct {
	ID          int      `json:"id" db:"id"`
	RoleName    string   `json:"name" db:"role_name"`
	Description string   `json:"description" db:"description"`
	RoleList    string   `json:"routes" db:"role_list"`
	Api         []string `json:"api"`
}

type TimingModel struct {
	Page   int  `json:"page"`
	Limit  int  `json:"limit"`
	Status *int `json:"status"`
	Action *int `json:"action"`
}

type DataxInfoListReq struct {
	Remark string `json:"remark"`
	Typ    string `json:"typ"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type DataxInfoInsertReq struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
	Remark   string `json:"remark"`
	Typ      string `json:"typ"`
}

type DataxInfoDelReq struct {
	ID int `json:"id"`
}

type DataxInfoTestLinkReq struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
	Remark   string `json:"remark"`
	Typ      string `json:"typ"`
}

func (this DataxInfoInsertReq) Validate() (err error) {
	if this.IP == "" {
		err = errors.New("ip 不能为空")
	}
	return
}

type TransferReq struct {
	AutoIncrementId string `json:"autoIncrementId"`
	EsConnect       int    `json:"es_connect"`
	SelectType      string `json:"selectType"`
	Remark          string `json:"remark"`
	SelectTable     string `json:"selectTable"`
	TypeName        string `json:"type_name"`
	Cols            struct {
		TableCols []string `json:"tableCols"`
		EsCols    []struct {
			Col   string `json:"col"`
			TbCol string `json:"tbCol"`
		} `json:"esCols"`
	} `json:"cols"`
	IndexName       string `json:"indexName"`
	Reset           bool   `json:"reset"`
	GoNum           int    `json:"goNum"`
	EsFlushInterval int    `json:"esFlushInterval"`
	EsBufferSize    int    `json:"esBufferSize"`
	EsDocID         string `json:"esDocId"`
	BufferSize      int    `json:"bufferSize"`
	MaxIdleConns    int    `json:"maxIdleConns"`
	MaxOpenConns    int    `json:"maxOpenConns"`
	CrontabSpec     string `json:"crontab_spec"`
}

type SelectType struct {
	ID     int    `json:"id"`
	Remark string `json:"remark"`
	Typ    string `json:"typ"`
}

type TransferLogReq struct {
	EsConnect int `json:"es_connect"`
}

func (this *TransferReq) ParseSelectType() (*SelectType, error) {
	selectType := new(SelectType)

	err := json.Unmarshal([]byte(this.SelectType), selectType)
	if err != nil {
		return nil, err
	}
	return selectType, nil
}
