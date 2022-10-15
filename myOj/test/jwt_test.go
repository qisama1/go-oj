package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// 生成token
func TestGenerateToken(t *testing.T) {
	UserClaim := &UserClaims{
		Identity:       "user_1",
		Name:           "Get",
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	signedString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}
	// 获取的token
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.4inO9HZINmKFYO9qEF2SYYPHk0GuuA-qUdwIhUa8USE
	fmt.Println(signedString)
}

// 解析token
func TestAnalyseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.4inO9HZINmKFYO9qEF2SYYPHk0GuuA-qUdwIhUa8USE"
	claim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(claims)
	if claims.Valid {
		fmt.Println(claim)
	}
}
