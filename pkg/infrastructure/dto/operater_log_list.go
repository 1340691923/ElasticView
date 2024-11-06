package dto

type GmOperaterLogList struct {
	Page           int      `json:"page"`
	Limit          int      `json:"limit"`
	UserId         int      `json:"operater_id"`
	OperaterAction string   `json:"operater_action"`
	Date           []string `json:"date"`
}
