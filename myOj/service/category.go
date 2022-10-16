package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myOj/define"
	"myOj/models"
	"myOj/utils"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags 内部方法
// @Summary 分类列表
// @Param authorization header string true "authorization"
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Param keyword query string false "查询关键词"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-list [get]
func GetCategoryList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		fmt.Println("Get ProblemList Size strconv err ", err)
		return
	}

	// page是1，其实是从0开始的offset
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage)) // 设置默认值
	if err != nil {
		fmt.Println("Get ProblemList Page strconv err ", err)
		return
	}
	page = (page - 1) * size // 起始位置
	var count int64
	keyword := c.Query("keyword")
	list := make([]models.CategoryBasic, 0)
	tx := models.DB.Model(new(models.CategoryBasic))
	if keyword != "" {
		tx.Where("name like %" + keyword + "%")
	}
	res := tx.Limit(size).Offset(page).Find(&list)
	err = res.Error
	count = res.RowsAffected
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "获取分类类别失败",
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

// CategoryCreate
// @Tags 内部方法
// @Summary 创建分类
// @Param Authorization header string true "Authorization"
// @Param name formData string true "name"
// @Param parent_id formData int false "parent_id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-create [post]
func CategoryCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分类名不能为空",
		})
		return
	}
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	category := models.CategoryBasic{
		Name:     name,
		Identity: utils.GetUUID(),
		ParentId: parentId,
	}
	err := models.DB.Create(&category).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Create err, " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建分类成功",
	})

}

// CategoryModify
// @Tags 内部方法
// @Summary 修改分类
// @Param Authorization header string true "Authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parent_id formData int false "parent_id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-modify [put]
func CategoryModify(c *gin.Context) {
	identity := c.PostForm("identity")
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	category := models.CategoryBasic{
		Name:     name,
		Identity: identity,
		ParentId: parentId,
	}
	err := models.DB.Where("identity = ?", identity).Updates(&category).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Modify err, " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "修改分类成功",
	})
}

// CategoryDelete
// @Tags 内部方法
// @Summary 删除分类
// @Param Authorization header string true "Authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-delete [delete]
func CategoryDelete(c *gin.Context) {
	identity := c.Query("identity")
	var cnt int64
	err := models.DB.Model(new(models.ProblemCategory)).Where("category_id = (SELECT id FROM category_basic WHERE identity = ? LIMIT 1)", identity).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Delete err, " + err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "表关联问题，不可以删除",
		})
		return
	}
	err = models.DB.Where("identity = ?", identity).Delete(new(models.CategoryBasic)).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Delete err, " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除分类成功",
	})
}
