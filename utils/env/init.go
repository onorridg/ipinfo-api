package env

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	api "api"
	middleware "middleware"
	password "password"
	redis "redisDB"
)

type ENV struct {
	API_PORT                   string
	JWT_SECRET_KEY             string
	PASSWORD_SALT              string
	REDIS_HOST                 string
	REDIS_PASSWORD             string
	REDIS_USER_DB_ID           int
	REDIS_CACHE_TIMEOUT_SECOND time.Duration
}

func InitEnvVars() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(filepath.Join(wd, ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rCacheTime, err := strconv.Atoi(os.Getenv("REDIS_CACHE_TIMEOUT_SECOND"))
	if err != nil {
		log.Fatal(err)
	}
	rUserDID, err := strconv.Atoi(os.Getenv("REDIS_USER_DB_ID"))
	if err != nil {
		log.Fatal(err)
	}

	env := ENV{
		API_PORT:                   os.Getenv("API_PORT"),
		JWT_SECRET_KEY:             os.Getenv("JWT_SECRET_KEY"),
		PASSWORD_SALT:              os.Getenv("PASSWORD_SALT"),
		REDIS_HOST:                 os.Getenv("REDIS_HOST"),
		REDIS_PASSWORD:             os.Getenv("REDIS_PASSWORD"),
		REDIS_USER_DB_ID:           rUserDID,
		REDIS_CACHE_TIMEOUT_SECOND: time.Duration(rCacheTime),
	}
	api.InitIPInfoVars(
		env.API_PORT,
		env.REDIS_HOST,
		env.REDIS_PASSWORD,
		env.REDIS_CACHE_TIMEOUT_SECOND,
	)
	redis.InitRedisVars(
		env.REDIS_HOST,
		env.REDIS_PASSWORD,
		env.REDIS_USER_DB_ID,
	)
	password.InitPasswordVars(
		env.PASSWORD_SALT,
	)
	middleware.InitAuthVars(
		env.JWT_SECRET_KEY,
	)
}
