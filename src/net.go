package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// 读取客户端发送的命令
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			return
		}

		// 处理命令
		cmd = strings.TrimSpace(cmd) // 去掉两边空格
		parts := strings.Split(cmd, " ")
		switch parts[0] {
		case "PING":
			// 响应PING命令
			if len(parts) > 1 {
				response := strings.Join(parts[1:], " ")
				writer.WriteString(fmt.Sprintf("PONG %s\r\n", response))
			} else {
				writer.WriteString("PONG\r\n")
			}
		case "SET":
			// 执行SET命令，保存键值对
			if len(parts) != 3 {
				writer.WriteString("ERR wrong number of arguments\r\n")
			} else {
				key := parts[1]
				value := parts[2]
				err := rdb.Set(key, value, 0).Err()
				if err != nil {
					writer.WriteString("ERR " + err.Error() + "\r\n")
				} else {
					writer.WriteString(fmt.Sprintf("OK SET %s %s\r\n", key, value))
				}
			}
		case "GET":
			// 执行GET命令，获取键对应的值
			if len(parts) != 2 {
				writer.WriteString("ERR wrong number of arguments\r\n")
			} else {
				key := parts[1]
				value, err := rdb.Get(key).Result()
				if err != nil {
					writer.WriteString("ERR " + err.Error() + "\r\n")
				} else {
					writer.WriteString(fmt.Sprintf("VALUE %s\r\n", value))
				}
			}
		case "ECHO":
			// 执行ECHO命令，返回命令参数
			if len(parts) > 1 {
				response := strings.Join(parts[1:], " ")
				writer.WriteString(fmt.Sprintf("%s\r\n", response))
			} else {
				writer.WriteString("\r\n")
			}
		case "QUIT":
			// 处理QUIT命令，关闭连接
			return
		default:
			writer.WriteString("ERR unknown command\r\n")
		}

		writer.Flush()
	}
}
