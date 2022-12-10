package ipinfo

//package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ip"
	"validator"
	m "middleware"
)

const (
	StatusUnknownError = 520
)

var authMiddleware = m.AuthMiddleware()

func getIpinfoV1(c *gin.Context) {
	IPStr := c.Param("ip")
	if !validator.IP(IPStr) {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message": "Invalid ip address"},
		)
		return
	}
	info, err := ip.Info(IPStr)
	if err != nil {
		c.IndentedJSON(
			StatusUnknownError,
			gin.H{"message": "Unknown error"},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, info)
}

func apiV1(router *gin.Engine) {
	V1 := router.Group("/v1")
	V1.Use(authMiddleware.MiddlewareFunc())
	{
		V1.GET("/:ip", getIpinfoV1)
	}
}

func login(router *gin.Engine){
	router.POST("/login", authMiddleware.LoginHandler)
}

func IPInfo() {
	router := gin.Default()
	login(router)
	apiV1(router)
	router.Run("localhost:8080")
}
