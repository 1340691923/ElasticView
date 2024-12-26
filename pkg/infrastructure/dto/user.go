package dto

type User struct {
	OAuthCode string `json:"oauth_code"`
	State     string `json:"state"`

	Username string `json:"username"`
	Password string `json:"password"`
}
