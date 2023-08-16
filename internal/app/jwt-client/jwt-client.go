package jwtclient

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClient struct {
	secret         []byte
	expTimeInHours int
	method         jwt.SigningMethod
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"userId"`
}

func New(secret string, expTimeInHours int) *JWTClient {
	return &JWTClient{
		secret:         []byte(secret),
		expTimeInHours: expTimeInHours,
		method:         jwt.SigningMethodHS256,
	}
}

func (c *JWTClient) GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(c.method, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix()},
		userId,
	})
	tokenString, err := token.SignedString(c.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (c *JWTClient) ParseToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return c.secret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}
