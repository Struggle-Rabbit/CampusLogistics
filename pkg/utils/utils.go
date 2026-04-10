package utils

import (
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
