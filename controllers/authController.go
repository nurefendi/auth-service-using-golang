package controllers

import (
	"auth-service/tools/helper"

	"github.com/gofiber/fiber/v2"
)

// TODO :
// - login
// - register
// - me
// - check access acl function

func AuthRegister(c *fiber.Ctx) {
	c.Status(fiber.StatusOK).
		JSON(helper.ResponseWithJson(nil))
}