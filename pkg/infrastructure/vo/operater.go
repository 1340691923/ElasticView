package vo

type GmOperaterLog struct {
	Id             int    `json:"id"`
	OperaterId     int    `json:"operater_id"`
	OperaterName   string `json:"operater_name"`
	OperaterAction string `json:"operater_action"`
	Method         string `json:"method"`
	Body           string `json:"body_str"`
	OperaterRoleId int    `json:"operater_role_id"`
	Created        string `json:"created"`
	CostTime       string `json:"cost_time"`
	Status         string `json:"status"`
}
