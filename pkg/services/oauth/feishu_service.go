package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type FeishuConfig struct {
	AppID       string `yaml:"appId"`
	AppSecret   string `yaml:"appSecret"`
	RedirectURI string `yaml:"redirectUri"`
}

type FeishuService struct {
	config *FeishuConfig
	log    *logger.AppLogger
}

func NewFeishuService(config *FeishuConfig, log *logger.AppLogger) *FeishuService {
	return &FeishuService{
		config: config,
		log:    log,
	}
}

func (s *FeishuService) GetAuthURL(state string) string {
	return fmt.Sprintf(
		"https://open.feishu.cn/open-apis/authen/v1/index?app_id=%s&redirect_uri=%s&state=%s",
		s.config.AppID,
		url.QueryEscape(s.config.RedirectURI),
		state,
	)
}

func (s *FeishuService) ExchangeToken(ctx context.Context, code string) (*TokenResponse, error) {
	tokenURL := "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"

	appTokenBody := map[string]string{
		"app_id":     s.config.AppID,
		"app_secret": s.config.AppSecret,
	}
	
	jsonAppTokenBody, err := json.Marshal(appTokenBody)
	if err != nil {
		s.log.Error("Failed to marshal Feishu app token request", zap.Error(err))
		return nil, err
	}

	appTokenReq, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(string(jsonAppTokenBody)))
	if err != nil {
		s.log.Error("Failed to create request for Feishu app token", zap.Error(err))
		return nil, err
	}
	
	appTokenReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	appTokenResp, err := client.Do(appTokenReq)
	if err != nil {
		s.log.Error("Failed to get Feishu app token", zap.Error(err))
		return nil, err
	}
	defer appTokenResp.Body.Close()

	appTokenBody, err = ioutil.ReadAll(appTokenResp.Body)
	if err != nil {
		s.log.Error("Failed to read Feishu app token response", zap.Error(err))
		return nil, err
	}

	var appTokenResult struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		AppAccessToken string `json:"app_access_token"`
		Expire int `json:"expire"`
	}

	if err := json.Unmarshal(appTokenBody, &appTokenResult); err != nil {
		s.log.Error("Failed to parse Feishu app token response", zap.Error(err), zap.String("body", string(appTokenBody)))
		return nil, err
	}

	if appTokenResult.Code != 0 {
		s.log.Error("Feishu app token error", zap.Int("code", appTokenResult.Code), zap.String("msg", appTokenResult.Msg))
		return nil, fmt.Errorf("Feishu error: %s (code: %d)", appTokenResult.Msg, appTokenResult.Code)
	}

	userTokenURL := "https://open.feishu.cn/open-apis/authen/v1/access_token"
	userTokenBody := map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	}
	
	jsonUserTokenBody, err := json.Marshal(userTokenBody)
	if err != nil {
		s.log.Error("Failed to marshal Feishu user token request", zap.Error(err))
		return nil, err
	}

	userTokenReq, err := http.NewRequestWithContext(ctx, "POST", userTokenURL, strings.NewReader(string(jsonUserTokenBody)))
	if err != nil {
		s.log.Error("Failed to create request for Feishu user token", zap.Error(err))
		return nil, err
	}
	
	userTokenReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	userTokenReq.Header.Set("Authorization", "Bearer "+appTokenResult.AppAccessToken)

	userTokenResp, err := client.Do(userTokenReq)
	if err != nil {
		s.log.Error("Failed to exchange Feishu user token", zap.Error(err))
		return nil, err
	}
	defer userTokenResp.Body.Close()

	userTokenRespBody, err := ioutil.ReadAll(userTokenResp.Body)
	if err != nil {
		s.log.Error("Failed to read Feishu user token response", zap.Error(err))
		return nil, err
	}

	var userTokenResult struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}

	if err := json.Unmarshal(userTokenRespBody, &userTokenResult); err != nil {
		s.log.Error("Failed to parse Feishu user token response", zap.Error(err), zap.String("body", string(userTokenRespBody)))
		return nil, err
	}

	if userTokenResult.Code != 0 {
		s.log.Error("Feishu user token error", zap.Int("code", userTokenResult.Code), zap.String("msg", userTokenResult.Msg))
		return nil, fmt.Errorf("Feishu error: %s (code: %d)", userTokenResult.Msg, userTokenResult.Code)
	}

	return &TokenResponse{
		AccessToken:  userTokenResult.Data.AccessToken,
		RefreshToken: userTokenResult.Data.RefreshToken,
		ExpiresIn:    userTokenResult.Data.ExpiresIn,
	}, nil
}

func (s *FeishuService) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	userInfoURL := "https://open.feishu.cn/open-apis/authen/v1/user_info"

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		s.log.Error("Failed to create request for Feishu user info", zap.Error(err))
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Failed to get Feishu user info", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Failed to read Feishu user info response", zap.Error(err))
		return nil, err
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			OpenID    string `json:"open_id"`
			UnionID   string `json:"union_id"`
			Email     string `json:"email"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		s.log.Error("Failed to parse Feishu user info response", zap.Error(err), zap.String("body", string(body)))
		return nil, err
	}

	if result.Code != 0 {
		s.log.Error("Feishu user info error", zap.Int("code", result.Code), zap.String("msg", result.Msg))
		return nil, fmt.Errorf("Feishu error: %s (code: %d)", result.Msg, result.Code)
	}

	return &UserInfo{
		ID:     result.Data.OpenID,
		Name:   result.Data.Name,
		Email:  result.Data.Email,
		Avatar: result.Data.AvatarURL,
	}, nil
}

func (s *FeishuService) GetProviderName() string {
	return "feishu"
}
