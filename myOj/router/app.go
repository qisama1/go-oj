package router

import (
	"github.com/gin-gonic/gin"
	"myOj/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 做些操作,路由规则
	r.GET("/ping", service.Ping)
	r.GET("/problem-list", service.GerProblemList)

	return r
}
