package jwt_test

import (
	jwtlocal "auth-service/middleware/jwt"
	"os"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
)

func TestGenerateAndParseToken(t *testing.T) {
	// require SECRET_KEY set for tests
	if os.Getenv("SECRET_KEY") == "" {
		t.Skip("SECRET_KEY not set")
	}

	claims := jwtlib.MapClaims{
		"userId": "00000000-0000-0000-0000-000000000000",
		"email":  "test@example.com",
		"exp":    time.Now().Add(5 * time.Minute).Unix(),
	}

	tokenStr, err := jwtlocal.GenerateToken(claims)
	if err != nil {
		t.Fatalf("generate token error: %v", err)
	}

	parsed, err := jwtlocal.JwtClaims(nil, tokenStr)
	if err != nil {
		t.Fatalf("parse token error: %v", err)
	}

	if parsed["email"] != "test@example.com" {
		t.Fatalf("unexpected claim email: %v", parsed["email"])
	}
}
