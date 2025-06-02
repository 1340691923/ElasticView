package dto

type UpdatePluginConfigReq struct {
	PluginID   string `json:"plugin_id" binding:"required"`
	AutoUpdate bool   `json:"auto_update"`
}

type GetPluginConfigReq struct {
	PluginID string `json:"plugin_id" binding:"required"`
}

type PluginConfigRes struct {
	PluginID   string `json:"plugin_id"`
	AutoUpdate bool   `json:"auto_update"`
}
