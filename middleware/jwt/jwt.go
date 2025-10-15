package jwt

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey []byte

func init() {
	// Load secret key and validate
	s := os.Getenv("SECRET_KEY")
	if s == "" {
		// In CI/Dev, allow a default non-production secret so tests can run.
		// In production, ensure SECRET_KEY is set.
		log.Warn("SECRET_KEY is not set; using default development secret")
		s = "dev-secret-CHANGE_ME"
	}
	secretKey = []byte(s)
}

func JwtClaims(c *fiber.Ctx, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		ip := ""
		if c != nil {
			ip = c.IP()
		}
		log.Error(ip, " ", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateToken returns a signed token string (not pointer) to simplify callers.
func GenerateToken(claim jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := t.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}
