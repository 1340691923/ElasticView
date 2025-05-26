package oauth

import (
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type OAuthInterface interface {
	GetOAuthUrl(callback string, state map[string]interface{}) string
	GetAppliactionName() string
	GetToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
	Enable() bool
	GetImg() string
	GetUserField() string
	GetConfig() map[string]interface{}
	SetConfig(data map[string]interface{})
}

func ProvideOAuthServiceRegistry(
	workWechat *WorkWechat,
	dingtalk *Dingtalk,
	feishu *Feishu,
) *OAuthServiceRegistry {
	return NewOAuthServiceRegistry(
		workWechat,
		dingtalk,
		feishu,
	)
}

type OAuthServiceRegistry struct {
	oAuthInterfaces []OAuthInterface
}

func NewOAuthServiceRegistry(services ...OAuthInterface) *OAuthServiceRegistry {
	return &OAuthServiceRegistry{services}
}

func (this *OAuthServiceRegistry) GetServices() []OAuthInterface {
	return this.oAuthInterfaces
}

func (this *OAuthServiceRegistry) FindServiceByName(name string) (OAuthInterface, error) {
	for _, v := range this.GetServices() {
		if v.GetAppliactionName() == name {
			return v, nil
		}
	}
	return nil, errors.New("没有找到对应第三方登录信息")
}
