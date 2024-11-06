package vo

type EsLink struct {
	ID               int             `json:"id"`
	Remark           string          `json:"remark"`
	Ip               string          `json:"ip"`
	Version          string          `json:"version"`
	CreateById       int             `json:"create_by_id"`
	CreateByUserName string          `json:"create_by_user_name"`
	Created          string          `json:"created"`
	Updated          string          `json:"updated"`
	EsLinkConfigs    []*EsLinkConfig `json:"es_link_configs"`
}

type EsLinkConfig struct {
	CfgRelationId int        `json:"cfg_relation_id"`
	Id            int        `json:"id"`
	Ip            string     `json:"ip"`
	Version       string     `json:"version"`
	EsLinkId      int        `json:"es_link_id"`
	User          string     `json:"user"`
	Pwd           string     `json:"pwd"`
	Remark        string     `json:"remark"`
	Created       string     `json:"created"`
	Updated       string     ` json:"updated" `
	RootPEM       string     ` json:"rootpem" `
	CertPEM       string     ` json:"certpem" `
	KeyPEM        string     `json:"keypem" `
	ShareRoles    []string   `json:"share_roles"`
	Header        []HeaderKv `json:"header"`
}

type EsLinkConfigV2 struct {
	Id         int        `json:"id"`
	User       string     `json:"user"`
	Pwd        string     `json:"pwd"`
	Remark     string     `json:"remark"`
	Created    string     `json:"created"`
	Updated    string     ` json:"updated" `
	RootPEM    string     ` json:"rootpem" `
	CertPEM    string     ` json:"certpem" `
	KeyPEM     string     `json:"keypem" `
	Header     []HeaderKv `json:"header"`
	ShareRoles []string   `json:"share_roles"`
}

type HeaderKv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EsLinkConfigOpt struct {
	Id int `json:"id"`

	Remark string `json:"remark"`
}

type EsLinkOpt struct {
	ID      int64  `json:"id"`
	Remark  string `json:"remark"`
	Version string `json:"version"`
}

type EsLinkTree struct {
	Label      string                `json:"label"`
	Value      string                `json:"value"`
	Selectable string                `json:"selectable"`
	Children   []*EsLinkTreeChildren `json:"children"`
}

func NewEsLinkTree(label string, value string) *EsLinkTree {
	return &EsLinkTree{Label: label, Value: value, Selectable: "false"}
}

type EsLinkTreeChildren struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	ParentId string `json:"parent_id"`
}

func NewEsLinkTreeChildren(label string, value string, parentId string) *EsLinkTreeChildren {
	return &EsLinkTreeChildren{Label: label, Value: value, ParentId: parentId}
}
