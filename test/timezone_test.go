package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTimezone(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	now := time.Now()
	fmt.Println("Local time:", now) // 2024-04-22 21:02:26.7572946 +0800 CST m=+0.633800601

	str := "2024-04-22 20:59:44"

	ti, _ := time.Parse(layout, str)
	fmt.Println("Parse time:", ti) // 2024-04-22 20:59:44 +0000 UTC

	loca, _ := time.LoadLocation("Asia/Shanghai")
	ti, _ = time.ParseInLocation(layout, str, loca)
	fmt.Println("Parse time with loca:", ti) // 2024-04-22 20:59:44 +0800 CST
}
