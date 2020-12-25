package jwtUtil

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	Email string
	UserID string
	CompanyID string
	Approved bool
}

func GenToken(user TokenClaims) (string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.UserID
	claims["company_id"] = user.CompanyID
	claims["approved"] = user.Approved
	claims["expiration"] = time.Now().Add(time.Hour*24).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Print("Token Generation Error")
		return "", err
	}
	return tokenString, nil
}

func CheckAndParseToken(tokenString string) (string, error){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		return claims["email"].(string), nil
	}
	return "", err
}