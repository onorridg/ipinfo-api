package ipinfo

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/onorridg/ipinfo-api/docs"

	"ip"
	m "middleware"
	p "password"
	rdb "redisDB"
	"validator"
)

var API_PORT string
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


// @Summary      getIPinfo
// @Security	 ApiKeyAuth
// @Description  get IP address info
// @Tags         ip info
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "IP Address"
// @Success      200  {object}  ip.IPData
// @Failure      400  {object}  string
// @Router       /ip/{ip} [get]
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
			v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

			account := v1.Group("/account")
			{
				account.POST("/sign-up", postSignInUserV1)
				account.POST("/sign-in", authMiddleware.LoginHandler)
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

func InitIPInfoVars(apiP, rH, rP string, rCacheTime time.Duration) {
	API_PORT = apiP
	REDIS_HOST = rH
	REDIS_PASSWORD = rP
	REDIS_CACHE_TIMEOUT_SECOND = rCacheTime
}


// @title           IP info API
// @version         1.0
// @description     Api для получения информации о IP адресе

// @contact.name   API Support
// @contact.url    https://t.me/onorridg

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func IPInfo() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routerHandler(router)

	router.Run(":" + API_PORT)
}
