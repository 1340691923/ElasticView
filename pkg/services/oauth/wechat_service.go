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
)

type WechatConfig struct {
	AppID       string `yaml:"appId"`
	AppSecret   string `yaml:"appSecret"`
	RedirectURI string `yaml:"redirectUri"`
}

type WechatService struct {
	config *WechatConfig
	log    *logger.AppLogger
}

func NewWechatService(config *WechatConfig, log *logger.AppLogger) *WechatService {
	return &WechatService{
		config: config,
		log:    log,
	}
}

func (s *WechatService) GetAuthURL(state string) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
		s.config.AppID,
		url.QueryEscape(s.config.RedirectURI),
		state,
	)
}

func (s *WechatService) ExchangeToken(ctx context.Context, code string) (*TokenResponse, error) {
	tokenURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		s.config.AppID,
		s.config.AppSecret,
		code,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", tokenURL, nil)
	if err != nil {
		s.log.Error("Failed to create request for WeChat token exchange", zap.Error(err))
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Failed to exchange WeChat token", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Failed to read WeChat token response", zap.Error(err))
		return nil, err
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		s.log.Error("Failed to parse WeChat token response", zap.Error(err), zap.String("body", string(body)))
		return nil, err
	}

	if result.ErrCode != 0 {
		s.log.Error("WeChat token exchange error", zap.Int("errcode", result.ErrCode), zap.String("errmsg", result.ErrMsg))
		return nil, fmt.Errorf("WeChat error: %s (code: %d)", result.ErrMsg, result.ErrCode)
	}

	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (s *WechatService) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	tokenInfoURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s",
		token,
		"", // We need to extract OpenID from the token or pass it separately
	)

	req, err := http.NewRequestWithContext(ctx, "GET", tokenInfoURL, nil)
	if err != nil {
		s.log.Error("Failed to create request for WeChat token validation", zap.Error(err))
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Failed to validate WeChat token", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Failed to read WeChat token validation response", zap.Error(err))
		return nil, err
	}

	var tokenInfo struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		OpenID  string `json:"openid"`
	}

	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		s.log.Error("Failed to parse WeChat token validation response", zap.Error(err), zap.String("body", string(body)))
		return nil, err
	}

	if tokenInfo.ErrCode != 0 {
		s.log.Error("WeChat token validation error", zap.Int("errcode", tokenInfo.ErrCode), zap.String("errmsg", tokenInfo.ErrMsg))
		return nil, fmt.Errorf("WeChat error: %s (code: %d)", tokenInfo.ErrMsg, tokenInfo.ErrCode)
	}

	userInfoURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",
		token,
		tokenInfo.OpenID,
	)

	req, err = http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		s.log.Error("Failed to create request for WeChat user info", zap.Error(err))
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		s.log.Error("Failed to get WeChat user info", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Failed to read WeChat user info response", zap.Error(err))
		return nil, err
	}

	var userInfo struct {
		OpenID     string   `json:"openid"`
		Nickname   string   `json:"nickname"`
		Sex        int      `json:"sex"`
		Province   string   `json:"province"`
		City       string   `json:"city"`
		Country    string   `json:"country"`
		HeadImgURL string   `json:"headimgurl"`
		Privilege  []string `json:"privilege"`
		UnionID    string   `json:"unionid"`
		ErrCode    int      `json:"errcode"`
		ErrMsg     string   `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		s.log.Error("Failed to parse WeChat user info response", zap.Error(err), zap.String("body", string(body)))
		return nil, err
	}

	if userInfo.ErrCode != 0 {
		s.log.Error("WeChat user info error", zap.Int("errcode", userInfo.ErrCode), zap.String("errmsg", userInfo.ErrMsg))
		return nil, fmt.Errorf("WeChat error: %s (code: %d)", userInfo.ErrMsg, userInfo.ErrCode)
	}

	return &UserInfo{
		ID:     userInfo.OpenID,
		Name:   userInfo.Nickname,
		Avatar: userInfo.HeadImgURL,
	}, nil
}

func (s *WechatService) GetProviderName() string {
	return "wechat"
}
