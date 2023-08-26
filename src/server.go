package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// 全局变量：用于键值对存储
var data = make(map[string]string)

func main() {
	// 监听 TCP 连接本地端口6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("监听失败！", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("正在监听端口: 6379")

	// 接收客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("连接失败: ", err.Error())
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// 读取客户端发送的命令
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取指令时出错: ", err.Error())
			return
		}

		// 处理命令
		result := processCommand(cmd)

		// 发送结果给客户端
		_, err = writer.WriteString(result + "\n")
		if err != nil {
			fmt.Println("响应错误！", err.Error())
			return
		}
		writer.Flush()
	}
}

func processCommand(cmd string) string {
	cmd = strings.TrimSuffix(cmd, "\n")
	// 解析命令和参数
	parts := strings.Split(cmd, " ")
	command := strings.ToUpper(parts[0])
	args := parts[1:]

	// 执行命令
	switch command {
	case "PING":
		return "PONG"
	case "ECHO":
		return strings.Join(args, " ")
	case "SET":
		if len(args) != 2 {
			return "SET 参数数量错误！"
		}
		key := args[0]
		value := args[1]
		data[key] = value
		return "OK"
	case "GET":
		if len(args) != 1 {
			return "GET 参数数量错误！"
		}
		// 获取键值对
		key := args[0]
		value, ok := data[key]
		if !ok {
			return "NIL"
		}
		return value
	case "DEL":
		if len(args) != 1 {
			return "DEL 参数数量错误！"
		}
		// 删除键值对
		key := args[0]
		delete(data, key)
		return "OK"
	case "QUIT":
		return "断开连接..."
	default:
		return fmt.Sprintf("未知指令: '%s'", command)
	}
}
