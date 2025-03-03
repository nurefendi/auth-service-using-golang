package enums

import "github.com/gofiber/fiber/v2"

type HttpMethod string

const (
	GET    HttpMethod = fiber.MethodGet
	PUT    HttpMethod = fiber.MethodPut
	DELETE HttpMethod = fiber.MethodDelete
	POST   HttpMethod = fiber.MethodPost
	PATCH  HttpMethod = fiber.MethodPatch
)

func (c HttpMethod) Name() string {
	switch c {
	case POST:
		return fiber.MethodPost
	case PUT:
		return fiber.MethodPut
	case DELETE:
		return fiber.MethodDelete
	case PATCH:
		return fiber.MethodPatch
	case GET:
	default:
		return fiber.MethodGet
	}
	return fiber.MethodGet
}

func GetValue(val string) HttpMethod {
	switch val {
	case fiber.MethodDelete:
		return DELETE
	case fiber.MethodPatch:
		return PATCH
	case fiber.MethodPost:
		return POST
	case fiber.MethodPut:
		return PUT
	case fiber.MethodGet:
	default:
		return GET
	}
	return GET
}