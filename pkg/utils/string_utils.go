package utils

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/db"
)

// RandomNumbers 返回指定长度的随机数字字符串
func RandomNumbers(length int) string {
	return randomStringWithBaseString("0123456789", length)
}

// RandomString 返回指定长度的随机字符的字符串
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

// AtoiOrDefault 将字符串 s 转化为 int，当 s 为空字符串时返回指定的默认值 defaultVal
func AtoiOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return val
}

// ParseInt64OrDefault 将字符串 s 转化为 int64，当 s 为空字符串时返回指定的默认值 defaultVal
func ParseInt64OrDefault(s string, defaultVal int64) int64 {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return val
}

// ParseInt64 将字符串 s 转换为 int64，解析错误抛出异常
func ParseInt64(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return val
}

// ParseFloatOrDefault 将字符串 s 转换为 float64，当 s 为空字符串时返回指定的默认值 defaultVal
func ParseFloatOrDefault(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return val
}

// NextId 生成 int64 类型的 id，专用于订单号生成
func NextId(ctx context.Context, keyPrefix string) int64 {
	// 1.生成时间戳
	now := time.Now().UTC()
	nowSecond := now.Unix()
	timestamp := nowSecond - constants.START_TIMESTAMP

	// 2.生成序列号
	date := now.Format("2006:01:02")
	// 2.2 自增长
	count, err := db.RedisClient.Incr(ctx, fmt.Sprintf("icr:%s:%s", keyPrefix, date)).Result()
	if err != nil {
		panic(err)
	}

	// 3.拼接并返回
	return int64((uint64(timestamp) << constants.COUNT_BITS) | uint64(count))
}
