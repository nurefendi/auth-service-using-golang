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
	// store pointer to avoid copying and to allow mutation by reference
	c.Locals(UserLocalKey, &local)
}

func GetLocals[T any](c *fiber.Ctx, key localKey) *T {
	value := c.Locals(key)
	// support both *T and T stored values
	switch v := value.(type) {
	case *T:
		return v
	case T:
		return &v
	default:
		_ = v
	}
	return nil
}
