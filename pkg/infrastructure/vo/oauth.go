package vo

type OAuthConfig struct {
	OauthUrl string `json:"oauthUrl"`
	Name     string `json:"name"`
	Enable   bool   `json:"enable"`
	Img      string `json:"img"`
}
