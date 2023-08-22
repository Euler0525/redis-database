package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

// 获取键对应的值
func getValue(key string) (string, error) {
	val, err := rdb.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("Key does not exist")
	} else if err != nil {
		return "", err
	}
	return val, nil
}
