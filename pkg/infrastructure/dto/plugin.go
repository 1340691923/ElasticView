package dto

type PluginMarketReq struct {
	SearchTxt   string `json:"search_txt"`
	OrderCol    string `json:"order_col"`
	OrderByDesc bool   `json:"order_by_desc"`
	IsInstall   *bool  `json:"is_install"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type InstallPlugin struct {
	PluginID string `json:"plugin_id"`
	Version  string `json:"version"`
}

type ImportEvKey struct {
	EvKey string `json:"ev_key"`
}
