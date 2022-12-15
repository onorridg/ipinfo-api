package auth

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	p "password"
	rdb "redisDB"
)

var JWT_SECRET_KEY string

var identityKey = "username"

type UserCredential struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserJWT struct {
	Code int `form:"code" json:"code" binding:"required"`
	Expire string `form:"expire" json:"expire" binding:"required"`
	Token string `form:"token" json:"token" binding:"required"`
}

var Roles = map[string][]string{
	"User": {"/v1/"},
}

type User struct {
	UserName string
	Role     []string
}

func InitAuthVars(jwtSK string) {
	JWT_SECRET_KEY = jwtSK
}


// @Summary      Sign-in
// @Description  sign-in
// @Tags         auth
// @Accept       multipart/form-data
// @Produce      json
// @Param        cUser		formData	UserCredential	true "User Credential"
// @Success      200  {object}  UserJWT
// @Failure      400  {object}  string
// @Router       /auth/sign-in [post]
func AuthMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(JWT_SECRET_KEY),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals UserCredential
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			ex, hPassword := rdb.GetValue(username)
			if ex == rdb.UserMissing {
				return nil, jwt.ErrFailedAuthentication
			} else if p.CompareHashPassword(hPassword, password) {
				return &User{
					UserName: username,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {

			// if v, ok := data.(*User); ok && (v.UserName == "admin" || v.UserName == "onorridg") {
			// 	return true
			// }
			//return false
			log.Println(c.Request.URL.Path)
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}
