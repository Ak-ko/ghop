package auth

import (
	"strconv"
	"time"

	"github.com/ak-ko/ghop.git/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.ENV.JWT_EXPIRATION)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID" : strconv.Itoa(userID),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}