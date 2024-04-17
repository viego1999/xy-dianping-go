package util

import "regexp"

func IsPhoneInvalid(phone string) bool {
	// 手机号码格式：11位数字，以1开头，第二位是3-9之间的数字，之后是其他9位数字
	return mismatch(phone, regexp.MustCompile(`^1[3-9]\d{9}$`))
}

func mismatch(str string, reg *regexp.Regexp) bool {
	return !reg.MatchString(str)
}
