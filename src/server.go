package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Database struct {
	data map[string]string
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]string),
	}
}

func (db *Database) Get(key string) string {
	return db.data[key]
}

func (db *Database) Set(key, value string) {
	db.data[key] = value
}

func (db *Database) Delete(key string) {
	delete(db.data, key)
}

type RedisServer struct {
	db *Database
}

func NewRedisServer() *RedisServer {
	return &RedisServer{
		db: NewDatabase(),
	}
}

func (s *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取指令时出错: %s\n", err)
			return
		}

		command = strings.TrimSpace(command)

		parts := strings.Split(command, " ")
		cmd := strings.ToLower(parts[0])
		args := parts[1:]

		switch cmd {
		case "ping":
			response := "+PONG\r\n"
			writer.WriteString(response)
		case "echo":
			if len(args) < 1 {
				response := "客户端ECHO指令格式错误！\r\n"
				writer.WriteString(response)
				continue
			}

			response := fmt.Sprintf("+%s\r\n", strings.Join(args, " "))
			writer.WriteString(response)
		case "get":
			if len(args) != 1 {
				response := "客户端GET指令格式错误！\r\n"
				writer.WriteString(response)
				continue
			}

			key := args[0]
			value := s.db.Get(key)

			if value != "" {
				response := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
				writer.WriteString(response)
			} else {
				response := "$-1\r\n"
				writer.WriteString(response)
			}
		case "set":
			if len(args) != 2 {
				response := "客户端SET指令格式错误！\r\n"
				writer.WriteString(response)
				continue
			}

			key := args[0]
			value := args[1]

			s.db.Set(key, value)

			response := "+OK\r\n"
			writer.WriteString(response)
		case "del":
			if len(args) < 1 {
				response := "客户端DEL指令格式错误！\r\n"
				writer.WriteString(response)
				continue
			}

			count := 0

			for _, key := range args {
				if _, ok := s.db.data[key]; ok {
					s.db.Delete(key)
					count++
				}
			}

			response := fmt.Sprintf(":%d\r\n", count)
			writer.WriteString(response)
		case "quit":
			fmt.Println("连接断开……")
			return
		default:
			response := "未知指令！\r\n"
			writer.WriteString(response)
		}

		writer.Flush()
	}
}

func main() {
	address := "localhost:6379"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("无法绑定到地址 %s: %s\n", address, err)
		return
	}

	fmt.Printf("正在监听地址 %s\n", address)

	server := NewRedisServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接时出错: %s\n", err)
			continue
		}

		go server.handleConnection(conn)
	}
}
