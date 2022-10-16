package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"myOj/define"
	"myOj/models"
	"myOj/utils"
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

// ProblemCreate
// @Tags 内部方法
// @Summary 创建问题
// @Param Authorization header string true "Authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData string true "max_runtime"
// @Param max_mem formData string true "max_mem"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array true "test_cases"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "title不能为空",
		})
		return
	}
	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "content不能为空",
		})
		return
	}
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))

	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))

	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")

	if len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "参数不能为空",
		})
		return
	}

	data := models.ProblemBasic{
		Identity:   utils.GetUUID(),
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}

	// 分类
	categoryBasic := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		categoryId, _ := strconv.Atoi(id)
		categoryBasic = append(categoryBasic, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(categoryId),
		})
	}
	data.ProblemCategories = categoryBasic

	// 测试用例 {"input":"1 2 \n" "output":"3"}
	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range testCases {
		// {"input":"1 2\n", "output": "3\n"}
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}

		testCaseBasic := &models.TestCase{
			Identity:  utils.GetUUID(),
			ProblemId: data.ID,
			Input:     caseMap["input"],
			Output:    caseMap["output"],
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}
	// 设置了外键以后居然可以自动绑定
	data.TestCases = testCaseBasics

	err := models.DB.Create(&data).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败, err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}
