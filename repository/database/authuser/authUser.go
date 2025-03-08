package authuser

import (
	"auth-service/config/database"
	"auth-service/dto"
	"auth-service/repository/dao"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func Save(c *fiber.Ctx) *fiber.Error {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data := locals.GetLocals[dao.AuthUser](c, locals.Entity)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return &fiber.Error{
			Code: fiber.StatusInternalServerError,
			Message: "Unable to conect database",
		}
	}
	saveRecord := db.Create(&data)
	if saveRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", saveRecord.Error.Error())
		return &fiber.Error{
			Code: fiber.StatusUnprocessableEntity,
			Message: saveRecord.Error.Error(),
		}
	}
	log.Info(currentAcess.RequestID, " Save data to DB, Affected rows :", saveRecord.RowsAffected)
	return nil
}

func FindByEmail(c *fiber.Ctx, email *string) (*dao.AuthUser, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return nil, &fiber.Error{
			Code: fiber.StatusInternalServerError,
			Message: "Unable to conect database",
		}
	}
	var data dao.AuthUser
	err := db.Model(&dao.AuthUser{}).
		Where("email = ?", email).
		Find(&data).Error
	if err != nil {
		log.Error(currentAcess.RequestID, err.Error())
		return nil, &fiber.Error{
			Code: fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return &data, nil
}

func FindById(c *fiber.Ctx, id uuid.UUID) (dao.AuthUser, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return dao.AuthUser{}, &fiber.Error{
			Code: fiber.StatusInternalServerError,
			Message: "Unable to conect database",
		}
	}
	var data dao.AuthUser
	err := db.Model(&dao.AuthUser{}).
		Where("id = ?", id).
		Find(&data).Error
	if err != nil {
		log.Error(currentAcess.RequestID, err.Error())
		return dao.AuthUser{}, &fiber.Error{
			Code: fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return data, nil
}

func FindByUserName(c *fiber.Ctx, username *string) (*dao.AuthUser, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return nil, &fiber.Error{
			Code: fiber.StatusInternalServerError,
			Message: "Unable to conect database",
		}
	}
	var data dao.AuthUser
	err := db.Model(&dao.AuthUser{}).
		Where("username = ?", username).
		Find(&data).Error
	if err != nil {
		log.Error(currentAcess.RequestID, err.Error())
		return nil, &fiber.Error{
			Code: fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return &data, nil
}