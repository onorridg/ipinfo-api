package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

const (
	UserMissing = 0
	UserExist = 1
)


func checkConnection(client *redis.Client)bool{
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}

func getConnection()*redis.Client{
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	if !checkConnection(rdb){
		log.Fatalln("Failed connection to redis database!")
	}
	return rdb
}

func SetKeyValue(key, value string)(error){
	rdb := getConnection()
	defer rdb.Close()

	result, err := rdb.SetNX(ctx, key, value, 0).Result()
	if err != nil {
		panic(err)
	} else if !result{
		return fmt.Errorf("Key already exist")
	}
	return nil
}

func GetValue(key string)(int, string){
	rdb := getConnection()
	defer rdb.Close()

	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return UserMissing , ""
	} else if err != nil{
		panic(err)
	}
	return UserExist, value
}