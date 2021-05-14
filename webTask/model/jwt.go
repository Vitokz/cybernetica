package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	jwt.StandardClaims
}
