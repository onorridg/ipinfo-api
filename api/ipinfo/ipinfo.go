package ipinfo

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ip"
	"validator"
)

const (
	StatusUnknownError = 520
)

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
	V1.GET("/:ip", getIpinfoV1)
}

func IPInfo() {
	router := gin.Default()
	apiV1(router)
	router.Run("localhost:8080")
}
