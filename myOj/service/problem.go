package service

import (
	"github.com/gin-gonic/gin"
	"myOj/models"
	"net/http"
)

func GerProblemList(c *gin.Context) {
	models.GetProblemList()
	c.String(http.StatusOK, "Get Problem List")
}
