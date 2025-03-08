package autrefreshtokens

import (
	"auth-service/config/database"
	"auth-service/dto"
	"auth-service/repository/dao"
	"auth-service/tools/locals"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func Save(c *fiber.Ctx, data *dao.AuthRefreshTokens) error {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return c.Status(fiber.StatusInternalServerError).
		SendString("Database error")
	}
	saveRecord := db.Save(&data)
	if saveRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", saveRecord.Error.Error())
		return saveRecord.Error
	}
	log.Info(currentAcess.RequestID, " Save data to DB, Affected rows :", saveRecord.RowsAffected)
	return nil
}

func DeleteByUserIdAndToken(c *fiber.Ctx, userId uuid.UUID, token string) error {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)

	if db == nil {
		log.Error(currentAccess.RequestID, " error: cannot find DB connection")
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	deleteRecord := db.Where("user_id = ? AND token = ?", userId, token).Delete(&dao.AuthRefreshTokens{})
	if deleteRecord.Error != nil {
		log.Error(currentAccess.RequestID, " error: ", deleteRecord.Error.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete record")
	}

	if deleteRecord.RowsAffected == 0 {
		log.Warn(currentAccess.RequestID, " warning: no record found to delete")
		return errors.New("failed to delete record")
	}

	log.Info(currentAccess.RequestID, " Deleted record, affected rows: ", deleteRecord.RowsAffected)
	return c.Status(fiber.StatusOK).SendString("Record deleted successfully")
}

func DeleteByUserId(c *fiber.Ctx, userId uuid.UUID) error {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)

	if db == nil {
		log.Error(currentAccess.RequestID, " error: cannot find DB connection")
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	deleteRecord := db.Where("user_id = ?", userId).Delete(&dao.AuthRefreshTokens{})
	if deleteRecord.Error != nil {
		log.Error(currentAccess.RequestID, " error: ", deleteRecord.Error.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete record")
	}

	if deleteRecord.RowsAffected == 0 {
		log.Warn(currentAccess.RequestID, " warning: no record found to delete")
		return errors.New("failed to delete record")
	}

	log.Info(currentAccess.RequestID, " Deleted record, affected rows: ", deleteRecord.RowsAffected)
	return c.Status(fiber.StatusOK).SendString("Record deleted successfully")
}


func FindByUserIdAndToken(c *fiber.Ctx, userId uuid.UUID, token string) (*dao.AuthRefreshTokens, error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data := dao.AuthRefreshTokens{}
	findRecord := db.Where("user_id = ? AND token = ? AND expires_at > CURRENT_TIMESTAMP", userId, token).
		First(&data)
	if findRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", findRecord.Error.Error())
		return nil, findRecord.Error
	}

	return &data, nil
}