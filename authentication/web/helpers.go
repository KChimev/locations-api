package main

import (
	"errors"
	"time"

	"github.com/kchimev/locations-api/authentication/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

func (a *app) generateToken() (string, error) {
	exp := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(constants.JWTKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *app) checkToken(tokenStr string) error {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWTKey), nil
	})

	if err != nil || !token.Valid {
		return errors.New(`cannot validate token`)
	}

	return nil
}
