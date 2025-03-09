package jwt

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func JwtClaims(c *fiber.Ctx, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		log.Error(c.IP(), " ", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GenerateToken(claim jwt.Claims) (*string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := jwt.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
	return &token, nil
}