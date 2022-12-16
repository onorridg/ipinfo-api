package ipinfo

import (
	"net/http"
	"time"

	"docs"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
var SWAGGER_DOCS_HOST string

const (
	StatusUnknownError = 520
)

type UserCredential struct {
	Username string `form:"username" json:"username" binding:"required" extensions:"x-order=0"`
	Password string `form:"password" json:"password" binding:"required" extensions:"x-order=1"`
}
type UserResetPassword struct {
	UserCredential
	NewPassword string `form:"new_password" json:"new_password" binding:"required" extensions:"x-order=3"`
}

type MessageResponse struct {
	Code    int    `form:"code" json:"code" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

var authMiddleware = m.AuthMiddleware()

//	@Summary		IP address info
//	@Security		ApiKeyAuth
//	@Description	get IP address info
//	@Tags			ip
//	@Accept			json
//	@Produce		json
//	@Param			ip	path		string	true	"IP Address"
//	@Success		200	{object}	ip.IPData
//	@Failure		400	{object}	MessageResponse
//	@Failure		401	{object}	MessageResponse
//	@Router			/ip/{ip} [get]
func getIpinfoV1(c *gin.Context) {
	IPStr := c.Param("ip")
	if !validator.IP(IPStr) {
		c.IndentedJSON(
			http.StatusBadRequest,
			MessageResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid ip address"},
		)
		return
	}
	info, err := ip.Info(IPStr)
	if err != nil {
		c.IndentedJSON(
			StatusUnknownError,
			MessageResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid ip address",
			},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, info)
}

//	@Summary		Sign-up
//	@Description	sign-up
//	@Tags			auth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			cUser	formData	UserCredential	true	"User Credential"
//	@Success		200		{object}	MessageResponse
//	@Failure		409		{object}	MessageResponse
//	@Failure		422		{object}	MessageResponse
//	@Router			/auth/sign-up [post]
func postSignUpUserV1(c *gin.Context) {
	uCred := UserCredential{}
	if err := c.ShouldBind(&uCred); err != nil {
		c.IndentedJSON(
			http.StatusUnprocessableEntity,
			MessageResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: "Invalid user credential",
			},
		)
		return
	}
	pHash := p.PasswordToHash(uCred.Password)
	err := rdb.SetKeyValue(uCred.Username, pHash)
	if err != nil {
		c.IndentedJSON(
			http.StatusConflict,
			MessageResponse{
				Code:    http.StatusConflict,
				Message: "Username already exist",
			},
		)
		return
	}
	c.IndentedJSON(
		http.StatusCreated,
		MessageResponse{
			Code:    http.StatusOK,
			Message: "User created",
		},
	)
}

//	@Summary		Password reset
//	@Description	password
//	@Tags			auth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			cUser	formData	UserResetPassword	true	"User Credential + new password"
//	@Success		200		{object}	MessageResponse
//	@Failure		409		{object}	MessageResponse
//	@Failure		422		{object}	MessageResponse
//	@Router			/auth/password [patch]
func patchPasswordResetV1(c *gin.Context) {
	uCred := UserResetPassword{}
	if err := c.ShouldBind(&uCred); err != nil {
		c.IndentedJSON(
			http.StatusUnprocessableEntity,
			MessageResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: "Invalid user credential",
			},
		)
		return
	}

	ex, hPassword := rdb.GetValue(uCred.Username)
	if ex == rdb.UserMissing || !p.CompareHashPassword(hPassword, uCred.Password) {
		c.IndentedJSON(
			http.StatusUnprocessableEntity,
			MessageResponse{
				Code:    http.StatusConflict,
				Message: "Username already exist",
			},
		)
		return
	}

	hNewPassword := p.PasswordToHash(uCred.NewPassword)
	_ = rdb.UpdateValue(uCred.Username, hNewPassword)
	c.IndentedJSON(
		http.StatusOK,
		MessageResponse{
			Code:    http.StatusOK,
			Message: "Password updated",
		},
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
			// Swagger Docs
			// base path swagger docs - /api/v1
			{
				docs.SwaggerInfo.Host = SWAGGER_DOCS_HOST
				router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
			}

			auth := v1.Group("/auth")
			{
				auth.POST("/sign-up", postSignUpUserV1)
				auth.POST("/sign-in", authMiddleware.LoginHandler)
				auth.PATCH("/password", patchPasswordResetV1)
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

func InitIPInfoVars(apiP, rH, rP, sDocs string, rCacheTime time.Duration) {
	API_PORT = apiP
	REDIS_HOST = rH
	REDIS_PASSWORD = rP
	REDIS_CACHE_TIMEOUT_SECOND = rCacheTime
	SWAGGER_DOCS_HOST = sDocs
}

func IPInfo() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routerHandler(router)

	router.Run(":" + API_PORT)
}
