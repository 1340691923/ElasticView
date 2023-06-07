// 节点配置层
package es_settings

import (
	"context"
	"encoding/json"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

type Settings struct {
	settingsV6 SettingsV6
	settingsV7 SettingsV7
	ver        int
}

// 节点设置
func NewSettingsByV6(client *elasticV6.Client) (*Settings, error) {
	settings := &Settings{}

	res, err := client.PerformRequest(context.Background(), elasticV6.PerformRequestOptions{
		Method: "GET",
		Path:   "/_cluster/settings?include_defaults=true",
	})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body, &settings.settingsV6)
	if err != nil {
		return nil, err
	}
	settings.ver = 6

	return settings, nil
}

func NewSettingsByV7(client *elasticV7.Client) (*Settings, error) {
	settings := &Settings{}

	res, err := client.PerformRequest(context.Background(), elasticV7.PerformRequestOptions{
		Method: "GET",
		Path:   "/_cluster/settings?include_defaults=true",
	})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body, &settings.settingsV7)
	if err != nil {
		return nil, err
	}

	settings.ver = 7

	return settings, nil
}

// 获取配置文件设置的 PathRepo
func (this *Settings) GetPathRepo() []string {
	if this.ver == 7 {
		list := []string{}
		for _, v := range this.settingsV7.Defaults.Path.Repo {
			list = append(list, v.(string))
		}
		return list
	}
	return this.settingsV6.Defaults.Path.Repo
}

// 获取配置文件中的 AllowedUrls
func (this *Settings) GetAllowedUrls() []interface{} {
	if this.ver == 7 {
		return this.settingsV7.Defaults.Repositories.URL.AllowedUrls
	}
	return this.settingsV6.Defaults.Repositories.URL.AllowedUrls
}
