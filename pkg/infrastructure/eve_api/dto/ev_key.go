package dto

type EvKeyReq struct {
	EvKey string `json:"ev_key"`
}

type Empty struct {
}

type FromEvPluginReq struct {
	SearchTxt          string   `json:"search_txt"`
	OrderByCol         string   `json:"order_by_col"`
	OrderByDesc        bool     `json:"order_by_desc"`
	Page               int      `json:"page"`
	Limit              int      `json:"limit"`
	HasDownloadPlugins []string `json:"has_download_plugins"`
	HasDownloadType    *bool    `json:"has_download_type"`
}

type StarPlugin struct {
	PluginId int64 `json:"plugin_id"`
}

type FormEvPluginInfoReq struct {
	PluginId int64 `json:"plugin_id"`
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
}

type GetPluginDownloadUrlReq struct {
	PluginAlias string `json:"plugin_alias"`
	Version     string `json:"version"`
	Os          string `json:"os"`
	Arch        string `json:"arch"`
}

type GetEvPluginMaxVersion struct {
	PluginAlias string `json:"plugin_alias"`
}

type GetEvPluginsMaxVersion struct {
	PluginIds []string `json:"plugin_ids"`
}
