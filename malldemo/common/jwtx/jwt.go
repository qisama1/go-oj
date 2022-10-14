package jwtx

import "github.com/golang-jwt/jwt/v4"

func GetToken(secretKey string, iat, seconds, uid int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["uid"] = uid
	token := jwt.New(jwt.SigningMethodES256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
