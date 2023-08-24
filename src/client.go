package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	// 连接到 Redis 服务器
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("无法连接到 Redis 服务器: %s\n", err)
		return
	}
	defer conn.Close()

	// 创建读写缓冲区
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 读取用户输入并发送命令到 Redis 服务器
	for {
		fmt.Print("Redis> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// 发送命令到 Redis 服务器
		_, err := writer.WriteString(input + "\n")
		if err != nil {
			fmt.Printf("发送命令时出错: %s\n", err)
			return
		}
		writer.Flush()

		// 读取并打印 Redis 服务器的响应
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取响应时出错: %s\n", err)
			return
		}
		fmt.Println(response)

		// 如果用户输入 "quit"，退出客户端程序
		if strings.ToLower(input) == "quit" {
			break
		}
	}
}
