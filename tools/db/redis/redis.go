package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"strconv"
)

var (
	client *redis.Client
)

func init() {
	config := viper.GetStringMapString("redis")
	db, _ := strconv.Atoi(config["db"])
	client = redis.NewClient(&redis.Options{
		Addr:     config["host"],
		Password: config["pwd"],
		DB:       db,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("redis连接错误:", err)
	}
	fmt.Println("初始化连接redis, PONG:", pong)
}

func NewClient() *redis.Client {
	return client
}
