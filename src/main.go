package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func main() {
	// 创建 Redis 客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 密码置空
		DB:       0,  // 使用默认数据库
	})

	// 测试连接
	pong, err := rdb.Ping().Result()
	if err != nil {
		fmt.Println("Connect failed!", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

	handleRequest("SET A B")
	handleRequest("GET A")
	handleRequest("DEL A")

	// 关闭 Redis 连接
	err = rdb.Close()
	if err != nil {
		fmt.Println("Failed to close connection:", err)
	}
}
