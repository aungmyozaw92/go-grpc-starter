package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secretKey")

type JWTClaim struct {
  UserID int `json:"user_id"`
  jwt.RegisteredClaims
}

func GenerateJWT(userID int) (string, error) {
  claims := &JWTClaim{
    UserID: userID,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (int, error) {
  claims := &JWTClaim{}
  token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
    return jwtKey, nil
  })
  if err != nil || !token.Valid {
    return 0, errors.New("invalid token")
  }
  return claims.UserID, nil
}