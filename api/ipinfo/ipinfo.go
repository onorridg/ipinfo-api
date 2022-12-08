package main

import (	
	"net/http"

	"github.com/gin-gonic/gin"

	info "ip"

)


func getIpinfoV1(c *gin.Context){
	ip := c.Param("ip")
	i, _ := info.Info(ip)
	c.IndentedJSON(http.StatusOK, i)
}

func apiV1(router *gin.Engine){
	V1 := router.Group("/v1")
	V1.GET("/:ip", getIpinfoV1)
}

func main(){
	router := gin.Default()
	apiV1(router)
	router.Run("localhost:8080")
}