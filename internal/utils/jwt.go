package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

// sign and get the complete encoded token as a string using the secret token
func GenerateToken(id int, username string, email string, isAdmin bool) (string, error) {
	// expires in 30 days
	expirationTime := time.Now().Add(time.Hour * 24 * 30).Unix()

	claims := &Claims{
		ID:       id,
		Username: username,
		Email:    email,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
