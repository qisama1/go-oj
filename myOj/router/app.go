package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"myOj/middlewares"
	"myOj/service"
)

import _ "github.com/swaggo/gin-swagger" // gin-swagger middleware
import _ "github.com/swaggo/files"       // swagger embed files
import _ "myOj/docs"

func Router() *gin.Engine {
	r := gin.Default()

	// Swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 做些操作,路由规则

	// 问题路由
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)

	// 用户路由
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/user-login", service.Login)
	r.POST("/user-register", service.Register)
	r.POST("/user-code", service.SentCodeToRedis)

	// 排行榜
	r.GET("/user-rankList", service.GetRankList)

	// 提交记录
	r.GET("/submit-list", service.GetSubmitList)

	// 管理私有方法
	authGroup := r.Group("/admin", middlewares.AuthAdminCheck())
	// 问题创建
	authGroup.POST("/problem-create", service.ProblemCreate)
	// 问题修改
	authGroup.PUT("/problem-modify", service.ProblemModify)

	// 分类列表
	authGroup.GET("/category-list", service.GetCategoryList)
	// 分类创建
	authGroup.POST("/category-create", service.CategoryCreate)
	// 分类修改
	authGroup.PUT("/category-modify", service.CategoryModify)
	// 分类删除
	authGroup.DELETE("/category-delete", service.CategoryDelete)

	// 用户私有方法
	authUser := r.Group("/user", middlewares.AuthUserCheck())
	// 代码提交
	authUser.POST("/submit", service.Submit)
	return r
}
