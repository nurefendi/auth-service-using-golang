package locals

import (
	"auth-service/dto"

	"github.com/gofiber/fiber/v2"
)

type localKey string

const UserLocalKey localKey = "userAccess"

func SetLocals(c *fiber.Ctx, local dto.UserLocals) {
	c.Locals(UserLocalKey, local)
}

func GetLocals(c *fiber.Ctx) *dto.UserLocals {
	value := c.Locals(UserLocalKey)
	if userLocals, ok := value.(dto.UserLocals); ok {
		return &userLocals
	}
	return nil
}