package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"myOj/define"
	"myOj/models"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Param keyword query string false "查询关键词"
// @Param category_identity query string false "category"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
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
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")

	tx := models.GetProblemList(keyword, categoryIdentity) // 拿到了查到的DB

	list := make([]*models.ProblemBasic, 0)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get ProblemList err ", err)
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK,
			gin.H{
				"code": "-1",
				"msg":  "问题唯一标识不能为空",
			})
		return
	}
	data := new(models.ProblemBasic)
	err := models.DB.Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Where("identity = ?", identity).
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "Get ProblemDetail Err" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
