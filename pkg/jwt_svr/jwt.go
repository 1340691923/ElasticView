// JSON WEB TOKEN 层
package jwt_svr

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Jwt struct {
	cfg       *config.Config
	log       *logger.AppLogger
	jwtSecret string
}

func NewJwt(cfg *config.Config, log *logger.AppLogger) *Jwt {
	obj := &Jwt{cfg: cfg, log: log}
	obj.jwtSecret = cfg.AppSecret
	return obj
}

type Claims struct {
	UserID   int32  `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	RoleId   int32  `json:"role_id"`
	jwt.RegisteredClaims
}

// 生成token
func (this *Jwt) CreateToken(gmUser model.GmUserModel) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(7 * 24 * time.Hour) //token过期时间24小时
	//加入保存信息
	claims := Claims{
		gmUser.ID,
		gmUser.Username,
		gmUser.Realname,
		gmUser.RoleId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "xiao-wen-long",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //jwt.SigningMethodHS256 加密规则
	/*NewWithClaims(method SigningMethod, claims Claims)，method对应着SigningMethodHMAC struct{}
	，其包含SigningMethodHS256、SigningMethodHS384、SigningMethodHS512三种crypto.Hash方案
	*/

	token, err := tokenClaims.SignedString(util.Str2bytes(this.cfg.AppSecret)) // SignedString  该方法内部生成签名字符串，再用于获取完整、已签名的token
	if err != nil {
		err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_TOKEN)
		return token, err
	}
	return token, err
}

// 解析token
func (this *Jwt) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) { //ParseWithClaims   用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
		return util.Str2bytes(this.cfg.AppSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { //Valid 验证基于时间的声明exp, iat, nbf，注意如果没有任何声明在令牌中，仍然会被认为是有效的。并且对于时区偏差没有计算方法
			return claims, nil
		}
	}
	return nil, err
}
