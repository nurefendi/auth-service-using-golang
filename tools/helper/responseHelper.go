package helper

import "github.com/gofiber/fiber/v2"

func ResponseWithJson(data interface{}) interface{} {
	return fiber.Map{
		"data":   data,
	}
}