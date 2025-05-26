package oauth

import "context"

type OAuthProvider interface {
	GetAuthURL(state string) string
	
	ExchangeToken(ctx context.Context, code string) (*TokenResponse, error)
	
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)
	
	GetProviderName() string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}
