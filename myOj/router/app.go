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
	// 问题创建
	r.POST("/problem-create", middlewares.AuthAdminCheck(), service.ProblemCreate)
	// 分类列表
	r.GET("/category-list", middlewares.AuthAdminCheck(), service.GetCategoryList)
	// 分类创建
	r.POST("/category-create", middlewares.AuthAdminCheck(), service.CategoryCreate)
	// 分类修改
	r.PUT("/category-modify", middlewares.AuthAdminCheck(), service.CategoryModify)
	// 分类删除
	r.DELETE("/category-delete", middlewares.AuthAdminCheck(), service.CategoryDelete)

	return r
}
