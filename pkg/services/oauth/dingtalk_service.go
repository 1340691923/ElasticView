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

type DingtalkConfig struct {
	AppID       string `yaml:"appId"`
	AppSecret   string `yaml:"appSecret"`
	RedirectURI string `yaml:"redirectUri"`
}

type DingtalkService struct {
	config *DingtalkConfig
	log    *logger.AppLogger
}

func NewDingtalkService(config *DingtalkConfig, log *logger.AppLogger) *DingtalkService {
	return &DingtalkService{
		config: config,
		log:    log,
	}
}

func (s *DingtalkService) GetAuthURL(state string) string {
	return fmt.Sprintf(
		"https://login.dingtalk.com/oauth2/auth?client_id=%s&response_type=code&scope=openid&state=%s&redirect_uri=%s",
		s.config.AppID,
		state,
		url.QueryEscape(s.config.RedirectURI),
	)
}

func (s *DingtalkService) ExchangeToken(ctx context.Context, code string) (*TokenResponse, error) {
	tokenURL := "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"

	requestBody := map[string]string{
		"clientId":     s.config.AppID,
		"clientSecret": s.config.AppSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		s.log.Error("Failed to marshal DingTalk token request", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		s.log.Error("Failed to create request for DingTalk token exchange", zap.Error(err))
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Failed to exchange DingTalk token", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Failed to read DingTalk token response", zap.Error(err))
		return nil, err
	}

	var result struct {
		AccessToken string `json:"accessToken"`
		ExpiresIn   int    `json:"expireIn"`
		RefreshToken string `json:"refreshToken"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		s.log.Error("Failed to parse DingTalk token response", zap.Error(err), zap.String("body", string(body)))
		return nil, err
	}

	if result.AccessToken == "" {
		s.log.Error("DingTalk token exchange error", zap.String("response", string(body)))
		return nil, fmt.Errorf("DingTalk error: failed to get access token")
	}

	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (s *DingtalkService) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	userInfoURL := "https://api.dingtalk.com/v1.0/contact/users/me"

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		s.log.Error("Failed to create request for DingTalk user info", zap.Error(err))
		return nil, err
	}
	
	req.Header.Set("x-acs-dingtalk-access-token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Failed to get DingTalk user info", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Failed to read DingTalk user info response", zap.Error(err))
		return nil, err
	}

	var userInfo struct {
		Nick      string `json:"nick"`
		AvatarUrl string `json:"avatarUrl"`
		OpenID    string `json:"openId"`
		UnionID   string `json:"unionId"`
		Email     string `json:"email"`
		Mobile    string `json:"mobile"`
		StateCode string `json:"stateCode"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		s.log.Error("Failed to parse DingTalk user info response", zap.Error(err), zap.String("body", string(body)))
		return nil, err
	}

	if userInfo.OpenID == "" {
		s.log.Error("DingTalk user info error", zap.String("response", string(body)))
		return nil, fmt.Errorf("DingTalk error: failed to get user info")
	}

	return &UserInfo{
		ID:     userInfo.OpenID,
		Name:   userInfo.Nick,
		Email:  userInfo.Email,
		Avatar: userInfo.AvatarUrl,
	}, nil
}

func (s *DingtalkService) GetProviderName() string {
	return "dingtalk"
}
