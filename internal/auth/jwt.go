package auth

import (
	"errors"
	"github.com/ChrisCrawford1/Command/internal/models"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"time"
)

func GenerateAccessToken(user models.User) string {
	secret := []byte(os.Getenv("JWT_SIGN"))

	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)

	claims["userId"] = user.UUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	signedAccessToken, err := accessToken.SignedString(secret)

	if err != nil {
		log.Fatal(err)
	}

	return signedAccessToken
}

func ValidateAccessToken(providedToken string) (bool, jwt.MapClaims, error) {
	token, err := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SIGN")), nil
	})

	if !token.Valid || err != nil {
		return false, nil, err
	} else {
		return true, token.Claims.(jwt.MapClaims), nil
	}
}
