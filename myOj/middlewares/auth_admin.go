package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myOj/utils"
	"net/http"
)

// AuthAdminCheck 验证是不是有管理权限
func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: check user is admin?
		auth := c.GetHeader("Authorization")
		fmt.Println(auth)
		userClaim, err := utils.AnalyseToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized err" + err.Error(),
			})
			return
		}
		if userClaim.IsAdmin != 1 {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized",
			})
			return
		}
		c.Next()
	}
}
