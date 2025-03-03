package controllers

import (
	"auth-service/dto"
	"auth-service/tools/helper"
	"auth-service/tools/locals"
	"auth-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


func AuthRegister(c *fiber.Ctx) error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var request dto.AuthUserRegisterRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error(currentAccess.RequestID, " invalid bind json payload ")
		c.Status(fiber.StatusBadRequest).
		JSON(fiber.NewError(fiber.StatusBadRequest, " invalid bind json payload"))
		return nil
	}

	if err := helper.ValidateStruct(&request); err != nil {
		log.Error(currentAccess.RequestID, " Error validation ", err.Error())
		c.Status(fiber.StatusUnprocessableEntity).
		JSON(fiber.NewError(fiber.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	c.Locals(locals.PayloadLocalKey, request)
	fibererr := usecase.AuthUSeCase().Register(c)
	if fibererr != nil {
		c.Status(fibererr.Code).SendString(fibererr.Message)
		return nil
	}
	return c.SendStatus(fiber.StatusCreated)
}

