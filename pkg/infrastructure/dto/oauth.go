package dto

type GetOAuthConfigReq struct {
	CallBack string `json:"call_back"`
}

type SaveOAuthConfigReq struct {
	ApplicationName string                 `json:"application_name"`
	Config          map[string]interface{} `json:"config"`
}
