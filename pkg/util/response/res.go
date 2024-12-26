package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func JsonApiErr(c *gin.Context, code int, msg string, err error) {
	c.JSON(code, map[string]interface{}{"code": code, "msg": fmt.Sprintf(msg, err)})
}
