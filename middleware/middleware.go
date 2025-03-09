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
	return func(c *fiber.Ctx) error {
		c.Request().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		locals.SetLocals(c, dto.UserLocals{
			RequestID:    getRequestId(c),
			LanguageCode: getLanguageCode(c),
			ChannelID:    getChannelId(c),
		})
		return c.Next()
	}
}

func SetMiddlewareAUTH() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userAccess, err := getUserAccess(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil
		}
		locals.SetLocals(c, dto.UserLocals{
			RequestID:    getRequestId(c),
			LanguageCode: getLanguageCode(c),
			ChannelID:    getChannelId(c),
			UserAccess:   userAccess,
		})
		if err := getAuthorizationFunction(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).
				SendString("unauthorized access")
		}
		return c.Next()
	}
}

func SetMiddlewareAuthNoAcl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userAccess, err := getUserAccess(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil
		}
		locals.SetLocals(c, dto.UserLocals{
			RequestID:    getRequestId(c),
			LanguageCode: getLanguageCode(c),
			ChannelID:    getChannelId(c),
			UserAccess:   userAccess,
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

func getUserAccess(c *fiber.Ctx) (*dto.CurrentUserAccess, error) {
	tokenString := c.Cookies("token")
	var userAccess dto.CurrentUserAccess
	var isTokenExpired bool

	if tokenString != "" {
		dataClaim, err := jwtMidleware.JwtClaims(c, tokenString)
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					isTokenExpired = true
				} else {
					c.Status(fiber.StatusUnauthorized).
						SendString("Invalid token")
					return nil, errors.New("invalid token")
				}
			} else {
				c.Status(fiber.StatusUnauthorized).
					SendString("Invalid token")
				return nil, errors.New("invalid token")
			}

		}
		getuserId, _ := uuid.Parse(dataClaim["userId"].(string))
		userAccess = dto.CurrentUserAccess{
			UserID:   getuserId,
			UserName: dataClaim["userName"].(string),
			Email:    dataClaim["email"].(string),
		}
	} else {
		isTokenExpired = true
	}

	if isTokenExpired {
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil, errors.New("invalid token")
		}
		refreshClaim, err := jwtMidleware.JwtClaims(c, tokenString)
		if err != nil {
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil, errors.New("invalid token")
		}
		dataRefreshToken, err := authRefreshTokenRepository.FindByUserIdAndToken(c, refreshClaim["userId"].(uuid.UUID), refreshToken)
		if err != nil {
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil, errors.New("invalid token")
		}
		if dataRefreshToken.ID == uuid.Nil {
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil, errors.New("invalid token")
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

	}

	return &userAccess, nil
}
