package controllers

import (
	"auth-service/dto"
	"auth-service/tools/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// TODO :
// - login
// - register
// - me
// - check access acl function

func AuthRegister(c *fiber.Ctx) {

	var request dto.AuthUserRegisterRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error("invalid bind json payload ")
		c.Status(fiber.StatusBadRequest).
		JSON(err)
		return
	}

	if err := helper.ValidateStruct(&request); err != nil {
		log.Error(" Error validation ", err.Error())
		c.Status(fiber.StatusUnprocessableEntity).
		JSON(err)
		return
	}

	// add usecase here

	c.Status(fiber.StatusOK)
}

