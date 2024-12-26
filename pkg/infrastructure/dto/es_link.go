package dto

type InsertEsLink struct {
	Ip      string `json:"ip"`
	Remark  string `json:"remark"`
	Version string `json:"version"`
	CfgIds  []int  `json:"cfgIds"`
}

type InsertEsLinkCfg struct {
	User       string     `json:"user"`
	Pwd        string     `json:"pwd"`
	Remark     string     `json:"remark"`
	RootPEM    string     ` json:"rootpem" `
	CertPEM    string     ` json:"certpem" `
	KeyPEM     string     `json:"keypem" `
	ShareRoles []string   `json:"share_roles"`
	Header     []HeaderKv `json:"header"`
}

type UpdateEsLinkCfg struct {
	Id         int        `json:"id"`
	User       string     `json:"user"`
	Pwd        string     `json:"pwd"`
	Remark     string     `json:"remark"`
	RootPEM    string     ` json:"rootpem" `
	CertPEM    string     ` json:"certpem" `
	KeyPEM     string     `json:"keypem" `
	ShareRoles []string   `json:"share_roles"`
	LinkId     int        `json:"linkId"`
	Header     []HeaderKv `json:"header"`
}

type HeaderKv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DeleteEsLinkCfg struct {
	Id int `json:"id"`
}

type UpdateEsLink struct {
	Id      int    `json:"id"`
	Ip      string `json:"ip"`
	Remark  string `json:"remark" `
	Version string `json:"version"`
	CfgIds  []int  `json:"cfgIds"`
}

type GetEsCfgRelation struct {
	ID int `json:"id"`
}

type DeleteEsLink struct {
	Id int `json:"id"`
}
