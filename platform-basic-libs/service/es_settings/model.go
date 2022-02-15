package es_settings

type SettingsV6 struct {
	Defaults struct {
		Path struct {
			Repo []string `json:"repo"`
		} `json:"path"`
		Repositories struct {
			URL struct {
				AllowedUrls []interface{} `json:"allowed_urls"`
			} `json:"url"`
		} `json:"repositories"`
	} `json:"defaults"`
}

type SettingsV7 struct {
	Defaults struct {
		Path struct {
			Repo []interface{} `json:"repo"`
		} `json:"path"`
		Repositories struct {
			URL struct {
				AllowedUrls []interface{} `json:"allowed_urls"`
			} `json:"url"`
		} `json:"repositories"`
	} `json:"defaults"`
}
