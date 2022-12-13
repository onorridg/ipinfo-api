package ipinfo

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"

	"ip"
	m "middleware"
	p "password"
	rdb "redisDB"
	"validator"
)

var API_HOST string
var REDIS_HOST string
var REDIS_PASSWORD string
var REDIS_CACHE_TIMEOUT_SECOND time.Duration

const (
	StatusUnknownError = 520
)

type UserCredential struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type UserResetPassword struct {
	UserCredential
	NewPassword string `form:"new_password" json:"new_password" binding:"required"`
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

func postSignInUserV1(c *gin.Context) {
	uCred := UserCredential{}
	if err := c.ShouldBind(&uCred); err != nil {
		c.IndentedJSON(
			http.StatusUnprocessableEntity,
			gin.H{"message": "Invalid user credential"},
		)
		return
	}
	pHash := p.PasswordToHash(uCred.Password)
	err := rdb.SetKeyValue(uCred.Username, pHash)
	if err != nil {
		c.IndentedJSON(
			http.StatusConflict,
			gin.H{"message": "Username already exist"},
		)
		return
	}
	c.IndentedJSON(
		http.StatusCreated,
		gin.H{"message": "Ok"},
	)
}

func patchPasswordResetV1(c *gin.Context) {
	uCred := UserResetPassword{}
	if err := c.ShouldBind(&uCred); err != nil {
		c.IndentedJSON(
			http.StatusUnprocessableEntity,
			gin.H{"message": "Invalid user credential"},
		)
		return
	}

	ex, hPassword := rdb.GetValue(uCred.Username)
	if ex == rdb.UserMissing || !p.CompareHashPassword(hPassword, uCred.Password) {
		c.IndentedJSON(
			http.StatusUnprocessableEntity,
			gin.H{"message": "Invalid user credential"},
		)
		return
	}

	hNewPassword := p.PasswordToHash(uCred.NewPassword)
	_ = rdb.UpdateValue(uCred.Username, hNewPassword)
	c.IndentedJSON(
		http.StatusOK,
		gin.H{"message": "Password updated"},
	)
}

func routerHandler(router *gin.Engine) {
	rCache := persistence.NewRedisCache(
		REDIS_HOST,
		REDIS_PASSWORD,
		time.Minute,
	)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			account := v1.Group("/account")
			{
				account.POST("/signup", postSignInUserV1)
				account.POST("/signin", authMiddleware.LoginHandler)
				account.PATCH("/password", patchPasswordResetV1)
			}

			ip := v1.Group("/ip")
			ip.Use(authMiddleware.MiddlewareFunc())
			{
				ip.GET("/:ip", cache.CachePage(
					rCache,
					REDIS_CACHE_TIMEOUT_SECOND*time.Second,
					getIpinfoV1),
				)
			}

			statistic := v1.Group("/statistic")
			statistic.Use(authMiddleware.MiddlewareFunc())
			{
				// statistic
			}
		}
	}
}

func InitIPInfoVars(apiH, rH, rP string, rCacheTime time.Duration) {
	API_HOST = apiH
	REDIS_HOST = rH
	REDIS_PASSWORD = rP
	REDIS_CACHE_TIMEOUT_SECOND = rCacheTime
}

func IPInfo() {
	router := gin.Default()
	routerHandler(router)

	router.Run(API_HOST)
}
