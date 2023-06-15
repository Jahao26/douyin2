package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var mysecret = []byte("jahao26")

type Claims struct {
	UserId int64
	jwt.RegisteredClaims
}

func GenToken(id int64) (string, error) {
	claims := Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mysecret)
	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return mysecret, nil
	})

	return token, claims, err
}
