package middleware

import (
	"errors"
	"github.com/nurefendi/auth-service-using-golang/common/constants"
	"github.com/nurefendi/auth-service-using-golang/dto"
	jwtMidleware "github.com/nurefendi/auth-service-using-golang/middleware/jwt"
	authPermissionRepository "github.com/nurefendi/auth-service-using-golang/repository/database/authpermission"
	authRefreshTokenRepository "github.com/nurefendi/auth-service-using-golang/repository/database/autrefreshtokens"
	"github.com/nurefendi/auth-service-using-golang/tools/locals"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
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

func SetMiddlewareAUTH(withacl bool) fiber.Handler {
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
		if withacl {
			if err := getAuthorizationFunction(c); err != nil {
				return c.Status(fiber.StatusUnauthorized).
					SendString("unauthorized access")
			}
		}
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
	if tokenString == "" {
		return handleExpiredToken(c)
	}

	dataClaim, err := jwtMidleware.JwtClaims(c, tokenString)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return handleExpiredToken(c)
		}
		ip := c.IP()
		log.Error(ip, " Invalid token: ", err.Error())
		return nil, errors.New("invalid token")
	}
	// safe claim extraction
	userIDStr, _ := dataClaim["userId"].(string)
	userName, _ := dataClaim["userName"].(string)
	email, _ := dataClaim["email"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	userAccess := dto.CurrentUserAccess{
		UserID:   userID,
		UserName: userName,
		Email:    email,
	}

	return &userAccess, nil
}

func handleExpiredToken(c *fiber.Ctx) (*dto.CurrentUserAccess, error) {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		log.Error(c.IP(), " Empty refresh token")
		return nil, errors.New("invalid token")
	}

	refreshClaim, err := jwtMidleware.JwtClaims(c, refreshToken)
	if err != nil {
		log.Error(c.IP(), " Error claim: ", err.Error())
		return nil, errors.New("invalid token")
	}

	userIDStr, _ := refreshClaim["userId"].(string)
	userName, _ := refreshClaim["userName"].(string)
	email, _ := refreshClaim["email"].(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	dataRefreshToken, err := authRefreshTokenRepository.FindByUserIdAndToken(c, userID, refreshToken)
	if err != nil || dataRefreshToken.ID == uuid.Nil {
		log.Error(c.IP(), " Invalid refresh token: ", err)
		return nil, errors.New("invalid token")
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	userAccess := dto.CurrentUserAccess{
		UserID:   userID,
		UserName: userName,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err := jwtMidleware.GenerateToken(userAccess)
	if err == nil && token != "" {
		sameSite := "Strict"
		secure := true
		if os.Getenv("APP_ENV") == "local" {
			sameSite = "Lax"
			secure = false
		}
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  expirationTime,
			Secure:   secure,
			HTTPOnly: true,
			SameSite: sameSite,
		})
	}

	return &userAccess, nil
}
