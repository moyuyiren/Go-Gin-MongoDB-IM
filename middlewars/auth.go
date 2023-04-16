package middlewars

import (
	"GoIm/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthChechk() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		userClaims, err := utils.AnalyseToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户认证不通过",
			})
			return
		}
		c.Set("user_Claims", userClaims)
		c.Next()
	}
}
