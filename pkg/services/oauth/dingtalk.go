package oauth

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
	"time"
)

type Dingtalk struct {
	cfg *config.Config
}

func NewDingtalk(cfg *config.Config) *Dingtalk {
	return &Dingtalk{
		cfg: cfg,
	}
}

func (this *Dingtalk) GetAppliactionName() string {
	return "钉钉认证"
}

func (this *Dingtalk) Enable() bool {
	return this.cfg.DingtalkEnable()
}

func (this *Dingtalk) GetImg() string {
	return "dingtalk"
}

func (this *Dingtalk) GetOAuthUrl(callback string, state map[string]interface{}) string {
	if state == nil {
		state = map[string]interface{}{}
	}

	js, _ := json.Marshal(state)

	stateString := string(js)

	return fmt.Sprintf("https://login.dingtalk.com/oauth2/auth?"+
		"response_type=code"+
		"&client_id=%s"+
		"&scope=openid"+
		"&state=%s"+
		"&redirect_uri=%s",
		this.cfg.DingtalkAppId(),
		stateString,
		callback,
	)
}

type DingtalkTokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (this *Dingtalk) MockGetToken(code string) (*oauth2.Token, error) {
	raw := make(map[string]interface{})
	raw["code"] = code
	token := &oauth2.Token{
		AccessToken: "123",
		Expiry:      time.Unix(time.Now().Unix()+int64(60), 0),
	}
	token = token.WithExtra(raw)

	return token, nil
}

func (this *Dingtalk) GetToken(code string) (*oauth2.Token, error) {

	tokenUrl := "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"
	
	requestBody := map[string]string{
		"clientId": this.cfg.DingtalkAppId(),
		"clientSecret": this.cfg.DingtalkAppSecret(),
		"code": code,
		"grantType": "authorization_code",
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	
	data, err := util.PostJSON(tokenUrl, jsonBody)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenResp := struct {
		AccessToken string `json:"accessToken"`
		ExpiresIn   int    `json:"expireIn"`
		RefreshToken string `json:"refreshToken"`
	}{}
	
	err = json.Unmarshal(data, &tokenResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	token := &oauth2.Token{
		AccessToken: tokenResp.AccessToken,
		Expiry:      time.Unix(time.Now().Unix()+int64(tokenResp.ExpiresIn), 0),
	}

	raw := make(map[string]interface{})
	raw["code"] = code
	token = token.WithExtra(raw)

	return token, nil
}

func (this *Dingtalk) GetUserField() string {
	return "dingtalk_id"
}

type DingtalkUserInfo struct {
	Nick      string `json:"nick"`
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatarUrl"`
}

func (this *Dingtalk) MockGetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	userInfo := UserInfo{
		Id:          "dingtalk123",
		Username:    "dingtalk_user",
		DisplayName: "DingTalk User",
		Email:       "dingtalk@example.com",
		AvatarUrl:   "",
	}
	return &userInfo, nil
}

func (this *Dingtalk) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {

	userInfoUrl := "https://api.dingtalk.com/v1.0/contact/users/me"
	
	headers := map[string]string{
		"x-acs-dingtalk-access-token": token.AccessToken,
	}
	
	data, err := util.GetURLWithHeaders(userInfoUrl, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResp := struct {
		Nick      string `json:"nick"`
		OpenId    string `json:"openId"`
		UnionId   string `json:"unionId"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatarUrl"`
	}{}
	
	err = json.Unmarshal(data, &userResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userInfo := UserInfo{
		Id:          userResp.OpenId,
		Username:    userResp.Nick,
		DisplayName: userResp.Nick,
		Email:       userResp.Email,
		AvatarUrl:   userResp.AvatarUrl,
		UnionId:     userResp.UnionId,
	}

	if userInfo.Id == "" {
		userInfo.Id = userInfo.Username
	}

	return &userInfo, nil
}

func (this *Dingtalk) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"appId":     this.cfg.DingtalkAppId(),
		"appSecret": this.cfg.DingtalkAppSecret(),
		"enable":    this.cfg.DingtalkEnable(),
		"rootUrl":   this.cfg.GetRootUrl(),
	}
}

func (this *Dingtalk) SetConfig(data map[string]interface{}) {
	this.cfg.
		SetDingtalkAppId(cast.ToString(data["appId"])).
		SetDingtalkAppSecret(cast.ToString(data["appSecret"])).
		SetDingtalkEnable(cast.ToBool(data["enable"])).
		SetRootUrl(cast.ToString(data["rootUrl"])).
		GetViperInstance().
		WriteConfig()
}
