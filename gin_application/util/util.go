package util

import (
	"golang.org/x/exp/rand"
	"time"
)

// RandomString 生成一个随机的字符串
func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(uint64(time.Now().Unix()))
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
