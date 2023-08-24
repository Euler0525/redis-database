package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// 数据存储结构
var data = make(map[string]string)
var mutex sync.RWMutex

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done() // 并发编程：新任务减一

	fmt.Printf("新连接来自 %s...\n", conn.RemoteAddr().String())

	// 创建读写缓冲区
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 处理连接
	for {
		// 读取客户端发送的命令
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取命令时出错: %s\n", err)
			return
		}

		// 去除命令中的换行符
		command = strings.TrimSpace(command)

		// 解析命令和参数
		parts := strings.Split(command, " ")
		cmd := strings.ToLower(parts[0]) // 请求指令
		args := parts[1:]                // 指令参数

		// 执行命令
		switch cmd {
		case "ping":
			response := "PONG\n"
			writer.WriteString(response)
			writer.Flush()
		case "echo":
			response := strings.Join(args, " ") + "\n"
			writer.WriteString(response)
			writer.Flush()
		case "set":
			if len(args) != 2 {
				response := "ERROR: wrong number of arguments for 'SET' command\n"
				writer.WriteString(response)
				writer.Flush()
				continue
			}

			key := args[0]
			value := args[1]

			mutex.Lock()
			data[key] = value
			mutex.Unlock()

			response := "OK\n"
			writer.WriteString(response)
			writer.Flush()
		case "get":
			if len(args) != 1 {
				response := "ERROR: wrong number of arguments for 'GET' command\n"
				writer.WriteString(response)
				writer.Flush()
				continue
			}

			key := args[0]

			mutex.RLock()
			value, ok := data[key]
			mutex.RUnlock()

			if ok {
				response := value + "\n"
				writer.WriteString(response)
			} else {
				response := "(nil)\n"
				writer.WriteString(response)
			}

			writer.Flush()
		case "quit":
			fmt.Printf("连接断开: %s\n", conn.RemoteAddr().String())
			conn.Close()
			return
		default:
			response := "ERROR: unknown command\n"
			writer.WriteString(response)
			writer.Flush()
		}
	}
}

func main() {
	// 监听地址和端口
	address := "localhost:6379"

	// 创建TCP监听器
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("无法绑定到地址 %s: %s\n", address, err)
		return
	}

	fmt.Printf("正在监听地址 %s\n", address)

	// 创建等待所有连接完成的等待组
	var wg sync.WaitGroup

	// 接受连接并处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接时出错: %s\n", err)
			continue
		}

		// 增加等待组计数
		wg.Add(1)

		// 在单独的goroutine中处理连接
		go handleConnection(conn, &wg)
	}

	// 等待所有连接完成
	wg.Wait()
}
