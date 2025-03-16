package controllers

import (
	"auth-service/dto"
	"auth-service/tools/helper"
	"auth-service/tools/locals"
	"auth-service/usecase"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func SaveGroup(c *fiber.Ctx) error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var request dto.GroupDto

	if err := c.BodyParser(&request); err != nil {
		log.Error(currentAccess.RequestID, " invalid bind json payload ")
		c.Status(fiber.StatusBadRequest).
			JSON(fiber.NewError(fiber.StatusBadRequest, " invalid bind json payload"))
		return err
	}

	if err := helper.ValidateStruct(&request); err != nil {
		log.Error(currentAccess.RequestID, " Error validation ", err.Error())
		c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.NewError(fiber.StatusUnprocessableEntity, err.Error()))
		return err
	}
	fibererr := usecase.GroupUSeCase().Save(c, &request)
	if fibererr != nil {
		c.Status(fibererr.Code).SendString(fibererr.Message)
		return errors.New(fibererr.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}
func UpdateGroup(c *fiber.Ctx) error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var request dto.GroupDto

	if err := c.BodyParser(&request); err != nil {
		log.Error(currentAccess.RequestID, " invalid bind json payload ")
		c.Status(fiber.StatusBadRequest).
			JSON(fiber.NewError(fiber.StatusBadRequest, " invalid bind json payload"))
		return err
	}

	if err := helper.ValidateStruct(&request); err != nil {
		log.Error(currentAccess.RequestID, " Error validation ", err.Error())
		c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.NewError(fiber.StatusUnprocessableEntity, err.Error()))
		return err
	}
	fibererr := usecase.GroupUSeCase().Update(c, &request)
	if fibererr != nil {
		c.Status(fibererr.Code).SendString(fibererr.Message)
		return errors.New(fibererr.Error())
	}
	return c.SendStatus(fiber.StatusAccepted)
}
func GetGroupById(c *fiber.Ctx) error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)

	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Info(currentAccess.RequestID, " Invalid Id")
		c.Status(fiber.StatusBadRequest).
			JSON(fiber.NewError(fiber.StatusBadRequest, " Invalid Id "))
		return err
	}

	data, fibererr := usecase.GroupUSeCase().FindById(c, uid)
	if fibererr != nil {
		c.Status(fibererr.Code).SendString(fibererr.Message)
		return errors.New(fibererr.Error())
	}
	return c.Status(fiber.StatusOK).
		JSON(&data)
}
func DeleteGroupById(c *fiber.Ctx) error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)

	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Info(currentAccess.RequestID, " Invalid Id")
		c.Status(fiber.StatusBadRequest).
			JSON(fiber.NewError(fiber.StatusBadRequest, " Invalid Id "))
		return err
	}

	fibererr := usecase.GroupUSeCase().Delete(c, uid)
	if fibererr != nil {
		c.Status(fibererr.Code).SendString(fibererr.Message)
		return errors.New(fibererr.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func GetGroup(c *fiber.Ctx) error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var request dto.GroupPagination

	if err := c.BodyParser(&request); err != nil {
		log.Error(currentAccess.RequestID, " invalid bind json payload ")
		c.Status(fiber.StatusBadRequest).
			JSON(fiber.NewError(fiber.StatusBadRequest, " invalid bind json payload"))
		return err
	}

	data, total, fibererr := usecase.GroupUSeCase().FindAll(c, request)
	if fibererr != nil {
		log.Error(currentAccess.RequestID, " ", fibererr.Message)
		c.Status(fibererr.Code).SendString(fibererr.Message)
		return errors.New(fibererr.Error())
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"data":   data,
			"limit":  request.Limit,
			"offset": request.Offset,
			"search": request.Search,
			"total":  total,
		})
}