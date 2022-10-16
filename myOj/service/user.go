package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"myOj/models"
	"myOj/utils"
	"net/http"
	"strconv"
	"strings"
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

	token, err := utils.GenerateToken(data.Identity, data.Name, data.IsAdmin)
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

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param mail formData string false "mail"
// @Param name formData string false "name"
// @Param password formData string false "password"
// @Param code formData string false "code"
// @Param phone formData string false "phone"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-register [post]
func Register(c *gin.Context) {
	mail := c.PostForm("mail")
	if mail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "邮箱不能为空",
		})
		return
	}
	// 邮箱判重
	var count int64
	err := models.DB.Where("mail = ?", mail).Model(new(models.UserBasic)).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "获取用户失败, err" + err.Error(),
		})
	}
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "该邮箱已被注册, err" + err.Error(),
		})
	}
	code := c.PostForm("code")
	if code == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "code不能为空",
		})
		return
	}
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "name不能为空",
		})
		return
	}
	password := c.PostForm("password")
	if password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "password不能为空",
		})
		return
	}
	phone := c.PostForm("phone")
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "phone不能为空",
		})
		return
	}
	// 查看验证码是否正确
	redisCode, err := utils.Get(mail)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "get code err " + err.Error(),
		})
		return
	}
	if redisCode != code {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "验证码不正确",
		})
		return
	}
	// 数据的插入
	data := &models.UserBasic{
		Identity: utils.GetUUID(),
		Name:     name,
		Password: utils.GetMd5(password),
		Phone:    phone,
		Mail:     mail,
	}
	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "注册失败, err: " + err.Error(),
		})
		return
	}
	// 生成token，直接登录, 调用rpc可以吗或者rest
	url := "http://127.0.0.1:8080/user-login"
	// 表单数据
	contentType := "application/json"
	body, err := json.Marshal(data)
	_, err = http.Post(url, contentType, strings.NewReader(string(body)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "请重新登录",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

// SentCodeToRedis
// @Tags 公共方法
// @Summary 发送验证码
// @Param mail formData string false "mail"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-code [post]
func SentCodeToRedis(c *gin.Context) {
	mail := c.PostForm("mail")
	if mail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "邮箱为空",
		})
		return
	}
	err := utils.Set(mail, utils.GetRandom())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送验证码失败 , err:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送验证码成功",
	})
}

// GetRankList
// @Tags 公共方法
// @Summary 获取排行榜
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-rankList [get]
func GetRankList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "page必须是整数类型",
		})
		return
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "size必须是整数类型",
		})
		return
	}
	offset := size * (page - 1)
	list := make([]models.UserBasic, 0)
	var count int64

	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("finish_problem_num DESC, submit_problem_num ASC").
		Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "select err, err : " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
