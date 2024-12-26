package jwt_svr

import (
	"errors"
)

var TokenExpiredErr = errors.New("token has invalid claims: token is expired")
