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

type Feishu struct {
	cfg *config.Config
}

func NewFeishu(cfg *config.Config) *Feishu {
	return &Feishu{
		cfg: cfg,
	}
}

func (this *Feishu) GetAppliactionName() string {
	return "飞书认证"
}

func (this *Feishu) Enable() bool {
	return this.cfg.FeishuEnable()
}

func (this *Feishu) GetImg() string {
	return "feishu"
}

func (this *Feishu) GetOAuthUrl(callback string, state map[string]interface{}) string {
	if state == nil {
		state = map[string]interface{}{}
	}

	js, _ := json.Marshal(state)

	stateString := string(js)

	return fmt.Sprintf("https://open.feishu.cn/open-apis/authen/v1/index?"+
		"app_id=%s"+
		"&redirect_uri=%s"+
		"&state=%s",
		this.cfg.FeishuAppId(),
		callback,
		stateString,
	)
}

type FeishuTokenResp struct {
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	RefreshToken  string `json:"refresh_token"`
	RefreshExpire int    `json:"refresh_expires_in"`
}

func (this *Feishu) MockGetToken(code string) (*oauth2.Token, error) {
	raw := make(map[string]interface{})
	raw["code"] = code
	token := &oauth2.Token{
		AccessToken: "123",
		Expiry:      time.Unix(time.Now().Unix()+int64(60), 0),
	}
	token = token.WithExtra(raw)

	return token, nil
}

func (this *Feishu) GetToken(code string) (*oauth2.Token, error) {

	tokenUrl := "https://open.feishu.cn/open-apis/authen/v1/access_token"
	
	requestBody := map[string]string{
		"app_id": this.cfg.FeishuAppId(),
		"app_secret": this.cfg.FeishuAppSecret(),
		"grant_type": "authorization_code",
		"code": code,
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
		Code         int    `json:"code"`
		Msg          string `json:"msg"`
		Data         struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}{}
	
	err = json.Unmarshal(data, &tokenResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if tokenResp.Code != 0 {
		return nil, errors.Errorf("Feishu API error: code=%d, msg=%s", tokenResp.Code, tokenResp.Msg)
	}

	token := &oauth2.Token{
		AccessToken: tokenResp.Data.AccessToken,
		Expiry:      time.Unix(time.Now().Unix()+int64(tokenResp.Data.ExpiresIn), 0),
	}

	raw := make(map[string]interface{})
	raw["code"] = code
	token = token.WithExtra(raw)

	return token, nil
}

func (this *Feishu) GetUserField() string {
	return "feishu_open_id"
}

type FeishuUserInfo struct {
	Name      string `json:"name"`
	OpenId    string `json:"open_id"`
	UnionId   string `json:"union_id"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

func (this *Feishu) MockGetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	userInfo := UserInfo{
		Id:          "feishu123",
		Username:    "feishu_user",
		DisplayName: "Feishu User",
		Email:       "feishu@example.com",
		AvatarUrl:   "",
	}
	return &userInfo, nil
}

func (this *Feishu) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {

	userInfoUrl := "https://open.feishu.cn/open-apis/authen/v1/user_info"
	
	headers := map[string]string{
		"Authorization": "Bearer " + token.AccessToken,
	}
	
	data, err := util.GetURLWithHeaders(userInfoUrl, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResp := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Name      string `json:"name"`
			OpenId    string `json:"open_id"`
			UnionId   string `json:"union_id"`
			Email     string `json:"email"`
			AvatarUrl string `json:"avatar_url"`
		} `json:"data"`
	}{}
	
	err = json.Unmarshal(data, &userResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if userResp.Code != 0 {
		return nil, errors.Errorf("Feishu API error: code=%d, msg=%s", userResp.Code, userResp.Msg)
	}

	userInfo := UserInfo{
		Id:          userResp.Data.OpenId,
		Username:    userResp.Data.Name,
		DisplayName: userResp.Data.Name,
		Email:       userResp.Data.Email,
		AvatarUrl:   userResp.Data.AvatarUrl,
		UnionId:     userResp.Data.UnionId,
	}

	if userInfo.Id == "" {
		userInfo.Id = userInfo.Username
	}

	return &userInfo, nil
}

func (this *Feishu) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"appId":     this.cfg.FeishuAppId(),
		"appSecret": this.cfg.FeishuAppSecret(),
		"enable":    this.cfg.FeishuEnable(),
		"rootUrl":   this.cfg.GetRootUrl(),
	}
}

func (this *Feishu) SetConfig(data map[string]interface{}) {
	this.cfg.
		SetFeishuAppId(cast.ToString(data["appId"])).
		SetFeishuAppSecret(cast.ToString(data["appSecret"])).
		SetFeishuEnable(cast.ToBool(data["enable"])).
		SetRootUrl(cast.ToString(data["rootUrl"])).
		GetViperInstance().
		WriteConfig()
}
