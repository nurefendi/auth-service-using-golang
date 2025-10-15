package controllers

import (
	"github.com/nurefendi/auth-service-using-golang/dto"
	"github.com/nurefendi/auth-service-using-golang/middleware"
	"github.com/nurefendi/auth-service-using-golang/repository/database/adaudit"
	"github.com/nurefendi/auth-service-using-golang/tools/helper"
	"github.com/nurefendi/auth-service-using-golang/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UnlockRequest struct {
	Email string `json:"email"`
}

func UnlockAccount(c *fiber.Ctx) error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var req UnlockRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error(currentAccess.RequestID, "invalid bind")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "invalid payload"))
	}
	if err := helper.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.NewError(fiber.StatusUnprocessableEntity, err.Error()))
	}

	middleware.ResetFailedAttempts(req.Email)
	ip := c.IP()
	ua := c.Get("User-Agent")
	if err := adaudit.Save(c, "admin_unlock", nil, &ip, &ua, &req.Email); err != nil {
		log.Error(currentAccess.RequestID, "failed to save admin unlock audit: ", err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}
