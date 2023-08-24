package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

func main() {
	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 如果有密码，填写密码
		DB:       0,  // 选择数据库，默认为 0
	})

	// 创建上下文对象
	ctx := context.Background()

	// 创建读写缓冲区
	reader := bufio.NewReader(os.Stdin)

	// 读取用户输入并发送命令到 Redis 服务器
	for {
		fmt.Print("Redis> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// 解析命令和参数
		parts := strings.Split(input, " ")
		cmd := strings.ToLower(parts[0]) // 请求指令
		args := parts[1:]                // 指令参数

		// 执行命令
		switch cmd {
		case "ping":
			pong, err := client.Ping(ctx).Result()
			if err != nil {
				fmt.Printf("执行指令时出错: %s\n", err)
			} else {
				fmt.Println(pong)
			}
		case "echo":
			response := strings.Join(args, " ")
			fmt.Println(response)
		case "set":
			if len(args) != 2 {
				fmt.Println("指令格式错误！")
				continue
			}

			err := client.Set(ctx, args[0], args[1], 0).Err()
			if err != nil {
				fmt.Printf("执行指令时出错: %s\n", err)
			} else {
				fmt.Println("OK")
			}
		case "get":
			if len(args) != 1 {
				fmt.Println("指令格式错误！")
				continue
			}

			value, err := client.Get(ctx, args[0]).Result()
			if err == redis.Nil {
				fmt.Println("(nil)")
			} else if err != nil {
				fmt.Printf("执行指令时出错: %s\n", err)
			} else {
				fmt.Println(value)
			}
		case "del":
			if len(args) != 1 {
				fmt.Println("指令格式错误！")
				continue
			}

			err := client.Del(ctx, args[0]).Err()
			if err != nil {
				fmt.Printf("执行命令时出错: %s\n", err)
			} else {
				fmt.Println("OK")
			}
		case "quit":
			fmt.Println("退出客户端程序！")
			return
		default:
			fmt.Println("未知指令！")
		}
	}
}
