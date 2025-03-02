package middleware

import (
	"auth-service/common/constants"
	"auth-service/dto"
	"auth-service/middleware/jwt"
	"auth-service/tools/locals"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SetMiddlewareAUTH() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Missing or invalid token")
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		dataClaim, err := jwt.JwtClaims(c, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
			SendString("Invalid token")
		}
		userAccess := dto.CurrentUserAccess{
			UserID: dataClaim["userId"].(uuid.UUID),
			UserName: dataClaim["userName"].(string),
			Email: dataClaim["email"].(string),
			GroupID: dataClaim["groupId"].(uuid.UUID),
		}
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
	c.Path()
	// usecase for get 
	// AuthFunction inner join AuthPermission where  AuthFunction.path and AuthPermission.groupId

	panic("unimplement")
}