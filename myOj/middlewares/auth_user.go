package middlewares

import (
	"github.com/gin-gonic/gin"
	"myOj/utils"
	"net/http"
)

// AuthUserCheck 验证是不是用户登录了
func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: check user is admin?
		auth := c.GetHeader("Authorization")
		userClaim, err := utils.AnalyseToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized err" + err.Error(),
			})
			return
		}
		if userClaim == nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized",
			})
			return
		}
		c.Set("user", userClaim)
		c.Next()
	}
}
