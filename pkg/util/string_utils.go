package util

import (
	"math/rand"
	"time"
)

func RandomNumbers(length int) string {
	return randomStringWithBaseString("0123456789", length)
}

func RandomString(length int) string {
	return randomStringWithBaseString("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", length)
}

// randomStringWithBaseString 生成一个指定字符集合的指定长度的随机字符串
func randomStringWithBaseString(baseString string, length int) string {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 创建一个字符切片，用于存储最终的随机字符串
	b := make([]byte, length)

	// 为字符集创建一个切片，以便随机选择字符
	chars := []byte(baseString)

	// 循环生成随机字符串
	for i := 0; i < length; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}

	// 将字节切片转换为字符串
	return string(b)
}
