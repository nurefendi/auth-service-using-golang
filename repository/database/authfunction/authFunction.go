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

func FindById(c *fiber.Ctx, id uuid.UUID) (*dao.AuthFunction, error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data := dao.AuthFunction{}
	findRecord := db.Where("id = ?", id).
		Preload("Lang").
		First(&data)
	if findRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", findRecord.Error.Error())
		return nil, findRecord.Error
	}

	return &data, nil
}
