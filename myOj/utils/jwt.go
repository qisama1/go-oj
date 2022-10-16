package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// GenerateToken 生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: time.Now() + 3000, // 设置过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	// 获取的token
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.4inO9HZINmKFYO9qEF2SYYPHk0GuuA-qUdwIhUa8USE
	return tokenString, nil
}

// AnalyseToken 解析token
func AnalyseToken(token string) (*UserClaims, error) {
	claim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return claim, nil
}
