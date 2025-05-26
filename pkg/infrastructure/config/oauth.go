package config

type OAuthConfig struct {
	Wechat   *WechatConfig   `yaml:"wechat"`
	Dingtalk *DingtalkConfig `yaml:"dingtalk"`
	Feishu   *FeishuConfig   `yaml:"feishu"`
}

type WechatConfig struct {
	AppID       string `yaml:"appId"`
	AppSecret   string `yaml:"appSecret"`
	RedirectURI string `yaml:"redirectUri"`
}

type DingtalkConfig struct {
	AppID       string `yaml:"appId"`
	AppSecret   string `yaml:"appSecret"`
	RedirectURI string `yaml:"redirectUri"`
}

type FeishuConfig struct {
	AppID       string `yaml:"appId"`
	AppSecret   string `yaml:"appSecret"`
	RedirectURI string `yaml:"redirectUri"`
}
