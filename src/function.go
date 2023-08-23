package main

import (
	"fmt"
	"strings"
)

// 处理命令
func handleRequest(cmd string) {
	cmd = strings.TrimSpace(cmd) // 去掉两边空格
	parts := strings.Split(cmd, " ")
	switch parts[0] {
	case "PING":
		// 响应PING命令
		if len(parts) > 1 {
			strings.Join(parts[1:], " ")
		}
	case "SET":
		// 执行SET命令，保存键值对
		if len(parts) != 3 {
			fmt.Println("ERR wrong number of arguments")
		} else {
			key := parts[1]
			value := parts[2]
			err := rdb.Set(key, value, 0).Err()
			if err != nil {
				fmt.Println("Set Failed!")
			}
			fmt.Println(cmd + " Success!")
		}
	case "GET":
		// 执行GET命令，获取键对应的值
		if len(parts) != 2 {
			fmt.Println("ERR wrong number of arguments")
		} else {
			key := parts[1]
			value, err := rdb.Get(key).Result()
			if err != nil {
				fmt.Println("GET Failed!")
			}
			fmt.Println(cmd + " Success! value=" + value)
		}
	case "DEL":
		// 执行DEL命令，删除键值对
		if len(parts) != 2 {
			fmt.Println("Err wrong number of arguments")
		} else {
			_, err := rdb.Del(parts[1]).Result()
			if err != nil {
				fmt.Println("DEL Failed!")
			}
			fmt.Println(cmd + " Success!")
		}
	case "ECHO":
		// 执行ECHO命令，返回命令参数
		if len(parts) > 1 {
			strings.Join(parts[1:], " ")
		}
	case "QUIT":
		return
	default:
	}

}
