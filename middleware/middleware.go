package middleware

import (
	"auth-service/common/constants"
	"auth-service/dto"
	jwtMidleware "auth-service/middleware/jwt"
	authPermissionRepository "auth-service/repository/database/authpermission"
	authRefreshTokenRepository "auth-service/repository/database/autrefreshtokens"
	"auth-service/tools/locals"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SetMiddlewareJSON() fiber.Handler {
	return func(c *fiber.Ctx) error  {
		c.Request().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		locals.SetLocals(c, dto.UserLocals{
			RequestID: getRequestId(c),
			LanguageCode: getLanguageCode(c),
			ChannelID: getChannelId(c),
		})
		return c.Next()
	}
}

func SetMiddlewareAUTH() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Missing or invalid token")
		}
		dataClaim, err := jwtMidleware.JwtClaims(c, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}

		dataRefreshToken, err := authRefreshTokenRepository.FindByUserIdAndToken(c, dataClaim["userId"].(uuid.UUID), refreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}
		if dataRefreshToken.ID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}
		userAccess := dto.CurrentUserAccess{
			UserID: dataClaim["userId"].(uuid.UUID),
			UserName: dataClaim["userName"].(string),
			Email: dataClaim["email"].(string),
			GroupIDs: dataClaim["groupIds"].([]uuid.UUID),
		}
		expirationTime := time.Now().Add(15 * time.Minute)
		userAccess.StandardClaims = jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}
		token, _ := jwtMidleware.GenerateToken(userAccess)
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    *token,
			Expires:  expirationTime,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Strict",
		})

		locals.SetLocals(c, dto.UserLocals{
			RequestID: getRequestId(c),
			LanguageCode: getLanguageCode(c),
			ChannelID: getChannelId(c),
			UserAccess: userAccess,
		})

		err = getAuthorizationFunction(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("unauthorized access")
		}

		return c.Next()
	}
}

func SetMiddlewareAuthNoAcl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Missing or invalid token")
		}
		dataClaim, err := jwtMidleware.JwtClaims(c, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}

		dataRefreshToken, err := authRefreshTokenRepository.FindByUserIdAndToken(c, dataClaim["userId"].(uuid.UUID), refreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}
		if dataRefreshToken.ID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}
		userAccess := dto.CurrentUserAccess{
			UserID: dataClaim["userId"].(uuid.UUID),
			UserName: dataClaim["userName"].(string),
			Email: dataClaim["email"].(string),
			GroupIDs: dataClaim["groupIds"].([]uuid.UUID),
		}
		expirationTime := time.Now().Add(15 * time.Minute)
		userAccess.StandardClaims = jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}
		token, _ := jwtMidleware.GenerateToken(userAccess)
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    *token,
			Expires:  expirationTime,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Strict",
		})

		locals.SetLocals(c, dto.UserLocals{
			RequestID: getRequestId(c),
			LanguageCode: getLanguageCode(c),
			ChannelID: getChannelId(c),
			UserAccess: userAccess,
		})
		return c.Next()
	}
}

func getRequestId(c *fiber.Ctx) string {
	if c.Get(fiber.HeaderXRequestID) != "" {
		return c.Get(fiber.HeaderXRequestID)
	}
	return uuid.New().String()
}

func getChannelId(c *fiber.Ctx) string {
	channelID := c.Get(constants.CHANNEL_ID)
	if channelID == "" {
		channelID = constants.CHANNEL_SYSTEM
	}
	return strings.ToUpper(channelID)
}

func getLanguageCode(c *fiber.Ctx) string {
	if c.Get(fiber.HeaderContentLanguage) != "" {
		lang := c.Get(fiber.HeaderContentLanguage)
		code := strings.Split(lang, "-")[0]
		if code == "id" || code == "in" {
			return "in"
		}
		return code
	}
	return "en"
}

func getAuthorizationFunction(c *fiber.Ctx) error {
	can, err := authPermissionRepository.FindByGroupIdAndPathAndMethod(c)
	if err != nil {
		return err
	}
	if !can {
		return errors.New("unauthorized access")
	}
	return nil
}