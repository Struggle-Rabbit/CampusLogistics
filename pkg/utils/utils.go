package utils

import "golang.org/x/crypto/bcrypt"

func HashedPasswordFunc(password string) (string, error) {
	// 密码 bcrypt 加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
