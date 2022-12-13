package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
)

const (
	UserMissing = 0
	UserExist   = 1
)

var ctx = context.Background()

var redisHost string
var password string
var dbID int

func InitRedisVars(h, p string, db int) {
	redisHost = h
	password = p
	db = dbID
}

func checkConnection(client *redis.Client) bool {
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}

func GetConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: password,
		DB:       dbID,
	})
	if !checkConnection(rdb) {
		log.Fatalln("Failed connection to redis database!")
	}
	return rdb
}

func SetKeyValue(key, value string) error {
	rdb := GetConnection()
	defer rdb.Close()

	result, err := rdb.SetNX(ctx, key, value, 0).Result()
	if err != nil {
		panic(err)
	} else if !result {
		return fmt.Errorf("Key already exist")
	}
	return nil
}

func GetValue(key string) (int, string) {
	rdb := GetConnection()
	defer rdb.Close()

	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return UserMissing, ""
	} else if err != nil {
		panic(err)
	}
	return UserExist, value
}

func UpdateValue(key, value string) int {
	rdb := GetConnection()
	defer rdb.Close()

	result, err := rdb.SetXX(ctx, key, value, 0).Result()
	if err != nil {
		panic(err)
	} else if !result {
		return UserMissing
	}
	return UserExist
}
