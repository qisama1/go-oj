package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myOj/define"
	"myOj/models"
	"net/http"
	"strconv"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Param user_identity query string false "user identity"
// @Param problem_identity query string false "problem identity"
// @Param status query int false "状态"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		fmt.Println("Get ProblemList Size strconv err ", err)
	}

	// page是1，其实是从0开始的offset
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage)) // 设置默认值
	if err != nil {
		fmt.Println("Get ProblemList Page strconv err ", err)
	}
	page = (page - 1) * size // 起始位置
	var count int64

	userIdentity := c.Query("user_identity")
	problemIdentity := c.Query("problem_identity")
	status, err := strconv.Atoi(c.Query("status"))

	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	submits := make([]*models.SubmitBasic, 0)
	data := make(map[string]interface{})
	err = tx.Count(&count).Offset(page).Limit(20).Find(&submits).Error
	if err != nil {
		c.JSON(500, gin.H{
			"code":    -1,
			"message": "GetSubmitList err, " + err.Error(),
		})
		return
	}
	data["submits"] = submits
	data["count"] = count
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
