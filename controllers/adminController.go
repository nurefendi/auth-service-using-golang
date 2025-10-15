package controllers

import (
    "auth-service/dto"
    "auth-service/middleware"
    "auth-service/repository/database/adaudit"
    "auth-service/tools/helper"
    "auth-service/tools/locals"

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
