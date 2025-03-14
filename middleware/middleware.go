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
	"github.com/gofiber/fiber/v2/log"
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
	if c.Get(fiber.HeaderAcceptLanguage) != "" {
		lang := c.Get(fiber.HeaderAcceptLanguage)
		code := strings.Split(lang, "-")[0]
		if code == "id" || code == "in" {
			return "id"
		}
		return code
	}
	return "en"
}

func getAuthorizationFunction(c *fiber.Ctx) error {
	can, err := authPermissionRepository.FindByGroupIdAndPathAndMethod(c, c.Route().Path, c.Method())
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
			log.Error(c.IP(), " Empty token : ", err.Error())
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					log.Error(c.IP(), " token expired ")
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
			log.Error(c.IP(), " Empty refresh token")
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil, errors.New("invalid token")
		}
		refreshClaim, err := jwtMidleware.JwtClaims(c, refreshToken)
		if err != nil {
			log.Error(c.IP(), " Error claim ", err.Error())
			c.Status(fiber.StatusUnauthorized).
				SendString("Invalid token")
			return nil, errors.New("invalid token")
		}
		getuserId, _ := uuid.Parse(refreshClaim["userId"].(string))
		dataRefreshToken, err := authRefreshTokenRepository.FindByUserIdAndToken(c, getuserId, refreshToken)
		if err != nil {
			log.Error(c.IP(), " error empty : ", err.Error())
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
		userAccess = dto.CurrentUserAccess{
			UserID:   getuserId,
			UserName: refreshClaim["userName"].(string),
			Email:    refreshClaim["email"].(string),
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
