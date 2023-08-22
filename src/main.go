package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

// 定义一个全局变量，指向Redis客户端
var rdb *redis.Client

func initClient() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 密码置空
		DB:       0,  // 使用默认数据库
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

// func main() {
// 	err := initClient()
// 	if err != nil {
// 		fmt.Println("Failed to initialize Redis client:", err)
// 		return
// 	}
// 	fmt.Println("Connected!")

// 	// 在这里可以使用 rdb 进行各种 Redis 操作
// 	// rdb.Get("key")、rdb.Set("key", "value") ...

// 	// 关闭 Redis 连接
// 	err = rdb.Close()
// 	if err != nil {
// 		fmt.Println("Failed to close Redis client:", err)
// 		return
// 	}
// }

func main() {
	err := initClient()
	if err != nil {
		fmt.Println("Failed to initialize Redis client:", err)
		return
	}
	defer rdb.Close()

	fmt.Println("Connected!")

	value, err := getValue("myKey")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value)
	}
}
