// JSON WEB TOKEN 层
package jwt

import (
	"time"

	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = util.Str2bytes("1340691923@qq.com")

type Claims struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	RoleId   int32  `json:"role_id"`
	jwt.StandardClaims
}

//生成token
func GenerateToken(gmUser model.GmUserModel) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) //token过期时间24小时
	//加入保存信息
	claims := Claims{
		gmUser.ID,
		gmUser.Username,
		gmUser.RoleId,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "xiao-wen-long",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //jwt.SigningMethodHS256 加密规则
	/*NewWithClaims(method SigningMethod, claims Claims)，method对应着SigningMethodHMAC struct{}
	，其包含SigningMethodHS256、SigningMethodHS384、SigningMethodHS512三种crypto.Hash方案
	*/
	token, err := tokenClaims.SignedString(jwtSecret) // SignedString  该方法内部生成签名字符串，再用于获取完整、已签名的token
	if err != nil {
		err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_TOKEN)
		return token, err
	}
	return token, err
}

//解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) { //ParseWithClaims   用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { //Valid 验证基于时间的声明exp, iat, nbf，注意如果没有任何声明在令牌中，仍然会被认为是有效的。并且对于时区偏差没有计算方法
			return claims, nil
		}
	}
	return nil, err
}
