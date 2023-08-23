package main

import (
	"fmt"
	"net"

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

	fmt.Println("Connected")

	return nil
}

func main() {
	err := initClient()
	if err != nil {
		fmt.Println("Failed to initialize Redis client:", err)
		return
	}

	// 监听本地端口6379
	listener, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer rdb.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		go handleConnection(conn)
	}

}
