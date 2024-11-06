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

type WorkWechat struct {
	clientID     string
	clientSecret string
	agentId      string
	enable       bool
	cfg          *config.Config
}

func NewWorkWechat(cfg *config.Config) *WorkWechat {

	return &WorkWechat{
		cfg:          cfg,
		clientID:     cfg.OAuth.WorkWechat.Corpid,
		clientSecret: cfg.OAuth.WorkWechat.Secert,
		agentId:      cfg.OAuth.WorkWechat.AgentId,
		enable:       cfg.OAuth.WorkWechat.Enable,
	}
}

func (this *WorkWechat) GetAppliactionName() string {
	return "企业微信认证(内部应用)"
}

func (this *WorkWechat) Enable() bool {
	return this.enable
}

func (this *WorkWechat) GetImg() string {
	return "work_wechat"
}

func (this *WorkWechat) GetOAuthUrl(callback string, state map[string]interface{}) string {
	if state == nil {
		state = map[string]interface{}{}
	}

	js, _ := json.Marshal(state)

	stateString := string(js)

	return fmt.Sprintf("https://login.work.weixin.qq.com/wwlogin/sso/login/?"+
		"login_type=CorpApp"+
		"&appid=%s"+
		"&agentid=%s"+
		"&redirect_uri=%s"+
		"&state=%s",
		this.clientID,
		this.agentId,
		callback,
		stateString,
	)
}

type WecomInterToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (this *WorkWechat) MockGetToken(code string) (*oauth2.Token, error) {
	raw := make(map[string]interface{})
	raw["code"] = code
	token := &oauth2.Token{
		AccessToken: "123",
		Expiry:      time.Unix(time.Now().Unix()+int64(60), 0),
	}
	token = token.WithExtra(raw)

	return token, nil
}

func (this *WorkWechat) GetToken(code string) (*oauth2.Token, error) {

	//return this.MockGetToken(code)

	pTokenParams := &struct {
		CorpId     string `json:"corpid"`
		Corpsecret string `json:"corpsecret"`
	}{this.clientID, this.clientSecret}
	data, err := util.GetURL(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", pTokenParams.CorpId, pTokenParams.Corpsecret))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pToken := &WecomInterToken{}
	err = json.Unmarshal(data, pToken)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if pToken.Errcode != 0 {
		return nil, errors.Errorf("pToken.Errcode = %d, pToken.Errmsg = %s", pToken.Errcode, pToken.Errmsg)
	}

	token := &oauth2.Token{
		AccessToken: pToken.AccessToken,
		Expiry:      time.Unix(time.Now().Unix()+int64(pToken.ExpiresIn), 0),
	}

	raw := make(map[string]interface{})
	raw["code"] = code
	token = token.WithExtra(raw)

	return token, nil
}

func (this *WorkWechat) GetUserField() string {
	return "work_wechat_uid"
}

type WecomInternalUserResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	UserId  string `json:"UserId"`
	OpenId  string `json:"OpenId"`
}

type WecomInternalUserInfo struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Avatar  string `json:"avatar"`
	OpenId  string `json:"open_userid"`
	UserId  string `json:"userid"`
}

func (this *WorkWechat) MockGetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	userInfo := UserInfo{
		Id:          "aaa",
		Username:    "xwl",
		DisplayName: "xwl",
		Email:       "1340691923@qq.com",
		AvatarUrl:   "",
	}
	return &userInfo, nil
}

func (this *WorkWechat) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	//return this.MockGetUserInfo(token)
	accessToken := token.AccessToken
	code := token.Extra("code").(string)
	data, err := util.GetURL(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s", accessToken, code))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResp := &WecomInternalUserResp{}
	err = json.Unmarshal(data, userResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if userResp.Errcode != 0 {
		return nil, errors.Errorf("userIdResp.Errcode = %d, userIdResp.Errmsg = %s", userResp.Errcode, userResp.Errmsg)
	}
	if userResp.OpenId != "" {
		return nil, errors.Errorf("not an internal user")
	}
	// Use userid and accesstoken to get user information
	data, err = util.GetURL(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s", accessToken, userResp.UserId))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	infoResp := &WecomInternalUserInfo{}
	err = json.Unmarshal(data, infoResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if infoResp.Errcode != 0 {
		return nil, errors.Errorf("userInfoResp.errcode = %d, userInfoResp.errmsg = %s", infoResp.Errcode, infoResp.Errmsg)
	}

	userInfo := UserInfo{
		Id:          infoResp.UserId,
		Username:    infoResp.Name,
		DisplayName: infoResp.Name,
		Email:       infoResp.Email,
		AvatarUrl:   infoResp.Avatar,
	}

	if userInfo.Id == "" {
		userInfo.Id = userInfo.Username
	}

	return &userInfo, nil
}

func (this *WorkWechat) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"corpid":  this.clientID,
		"agentId": this.agentId,
		"secert":  this.clientSecret,
		"enable":  this.enable,
	}
}

func (this *WorkWechat) SetConfig(data map[string]interface{}) {
	this.clientID = cast.ToString(data["corpid"])
	this.agentId = cast.ToString(data["agentId"])
	this.clientSecret = cast.ToString(data["secert"])
	this.enable = cast.ToBool(data["enable"])
	this.cfg.
		SetWorkWechatAgentId(this.agentId).
		SetWorkWechatCorpid(this.clientID).
		SetWorkWechatSecert(this.clientSecret).
		SetWorkWechatEnable(this.enable).
		GetViperInstance().
		WriteConfig()
}
