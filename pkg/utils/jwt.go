package utils

import (
	"errors"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 JWT 声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleIDs  []uint `json:"role_ids"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 Access Token 和 Refresh Token
func GenerateToken(userID uint, username string, roleIDs []uint) (accessToken, refreshToken string, err error) {
	cfg := config.GlobalConfig.JWT

	// Access Token
	accessClaims := CustomClaims{
		UserID:   userID,
		Username: username,
		RoleIDs:  roleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.AccessExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "campus-logistics",
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return
	}

	// Refresh Token
	refreshClaims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.RefreshExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "campus-logistics",
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.Secret))
	return
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	cfg := config.GlobalConfig.JWT
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New(constant.MsgInvalidToken)
	}
	return claims, nil
}
