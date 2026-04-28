package utils

import (
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashedPasswordFunc 密码 bcrypt 加密
func HashedPasswordFunc(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPasswordFunc 加密校验
func VerifyPasswordFunc(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

// IsEmpty 判断是否为 nil 或空字符串
func IsEmpty(s *string) bool {
	return s == nil || *s == ""
}

// IsBlank 判断是否为 nil、空字符串或全空白
func IsBlank(s *string) bool {
	if s == nil {
		return true
	}
	return strings.TrimSpace(*s) == ""
}

// StrToInt 字符串转换为整数
func StrToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// StrToInt64 字符串转换为 int64 类型
func StrToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
