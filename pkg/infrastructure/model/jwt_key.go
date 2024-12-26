package model

type JwtKeyModel struct {
	JwtKey string
}

func (g *JwtKeyModel) TableName() string {
	return "jwt_key"
}
