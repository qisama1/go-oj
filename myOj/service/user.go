package service

import (
	"github.com/gin-gonic/gin"
	"myOj/models"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	userIdentity := c.Query("identity")
	if userIdentity == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Identity不能为空",
		})
		return
	}
	user := models.UserBasic{}
	err := models.DB.Omit("password").Where("identity = ?", userIdentity).First(&user).Error
	if err != nil {
		c.JSON(500, gin.H{
			"code":    -1,
			"message": "GetUserDetail err, " + userIdentity + " " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": user,
	})
}
