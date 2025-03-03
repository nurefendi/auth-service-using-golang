package locals

import (
	"auth-service/dto"

	"github.com/gofiber/fiber/v2"
)

type localKey string

const (
	UserLocalKey    localKey = "userAccess"
	PayloadLocalKey localKey = "payload"
	Entity          localKey = "entity"
)

func SetLocals(c *fiber.Ctx, local dto.UserLocals) {
	c.Locals(UserLocalKey, local)
}


func GetLocals[T any](c *fiber.Ctx, key localKey) *T {
	value := c.Locals(key)
	if locals, ok := value.(T); ok {
		return &locals
	}
	return nil
}
