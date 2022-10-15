package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"myOj/models"
	"myOj/utils"
	"net/http"
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

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}

	// md5
	password = utils.GetMd5(password)
	fmt.Println(password)
	data := new(models.UserBasic)
	err := models.DB.Where("name = ? AND password = ? ", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或者密码出错",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Login err" + err.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(data.Identity, data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "Generate token err " + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  data,
		"token": token,
	})
}
