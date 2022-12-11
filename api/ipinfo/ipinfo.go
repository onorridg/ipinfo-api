package ipinfo

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

type UserCredential struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

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

func postSignInUser(c *gin.Context){
	uCred := UserCredential{}
	if err := c.ShouldBind(&uCred); err != nil {
		c.IndentedJSON(
			http.StatusUnprocessableEntity,
			gin.H{"message": "Invalid user credential"},
		)
		return
	}
	// if username exist => 409 
}

func IPInfo() {
	router := gin.Default()
	router.POST("/signup", postSignInUser)
	router.POST("/signin", authMiddleware.LoginHandler)
	apiV1(router)
	router.Run("localhost:8080")
}
