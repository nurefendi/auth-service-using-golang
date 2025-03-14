package authfunction

import (
	"auth-service/config/database"
	"auth-service/dto"
	"auth-service/repository/dao"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func Save(c *fiber.Ctx, data *dao.AuthFunction) error {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return c.Status(fiber.StatusInternalServerError).
		SendString("Database error")
	}
	saveRecord := db.Save(data)
	if saveRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", saveRecord.Error.Error())
		return saveRecord.Error
	}
	log.Info(currentAcess.RequestID, " Save data to DB, Affected rows :", saveRecord.RowsAffected)
	return nil
}

func FindById(c *fiber.Ctx, id uuid.UUID) (*dao.AuthFunction, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data := dao.AuthFunction{}
	findRecord := db.Where("id = ?", id).
		Preload("Lang").
		First(&data)
	if findRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", findRecord.Error.Error())
		return nil, fiber.NewError(fiber.StatusInternalServerError, findRecord.Error.Error())
	}

	return &data, nil
}
func Delete(c *fiber.Ctx, portalID uuid.UUID) *fiber.Error {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)

	if db == nil {
		log.Error(currentAccess.RequestID, " error cannot find db connection")
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Unable to connect to database",
		}
	}

	deleteRecord := db.Delete(&dao.AuthFunction{}, "id = ?", portalID)
	if deleteRecord.Error != nil {
		log.Error(currentAccess.RequestID, " error ", deleteRecord.Error.Error())
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: deleteRecord.Error.Error(),
		}
	}

	log.Info(currentAccess.RequestID, " Deleted portal from DB, Affected rows: ", deleteRecord.RowsAffected)
	return nil
}