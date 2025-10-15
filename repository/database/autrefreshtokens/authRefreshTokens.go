package autrefreshtokens

import (
	"errors"
	"github.com/nurefendi/auth-service-using-golang/config/database"
	"github.com/nurefendi/auth-service-using-golang/dto"
	"github.com/nurefendi/auth-service-using-golang/repository/dao"
	"github.com/nurefendi/auth-service-using-golang/tools/helper"
	"github.com/nurefendi/auth-service-using-golang/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

// Save stores a refresh token record. Expectation: caller provides hashed token in data.Token.
func Save(c *fiber.Ctx, data *dao.AuthRefreshTokens) error {
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

// DeleteByUserIdAndToken finds a hashed token for the user matching the provided raw token and deletes it.
func DeleteByUserIdAndToken(c *fiber.Ctx, userId uuid.UUID, token string) error {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)

	if db == nil {
		log.Error(currentAccess.RequestID, " error: cannot find DB connection")
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	var candidates []dao.AuthRefreshTokens
	if err := db.Where("user_id = ? AND expires_at > CURRENT_TIMESTAMP", userId).Find(&candidates).Error; err != nil {
		log.Error(currentAccess.RequestID, " error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to find record")
	}

	for _, rec := range candidates {
		if helper.CompareHashBcript(token, rec.Token) == nil {
			// delete this record
			del := db.Delete(&dao.AuthRefreshTokens{}, "id = ?", rec.ID)
			if del.Error != nil {
				log.Error(currentAccess.RequestID, " error deleting: ", del.Error.Error())
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete record")
			}
			log.Info(currentAccess.RequestID, " Deleted record, affected rows: ", del.RowsAffected)
			return c.Status(fiber.StatusOK).SendString("Ok")
		}
	}

	log.Warn(currentAccess.RequestID, " warning: no record found to delete")
	return errors.New("failed to delete record")
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

// FindByUserIdAndToken searches for hashed refresh tokens for the user and verifies the provided raw token.
func FindByUserIdAndToken(c *fiber.Ctx, userId uuid.UUID, token string) (*dao.AuthRefreshTokens, error) {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var candidates []dao.AuthRefreshTokens
	if db == nil {
		log.Error(currentAccess.RequestID, " error cannot find db connection")
		return nil, c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	if err := db.Where("user_id = ? AND expires_at > CURRENT_TIMESTAMP", userId).Find(&candidates).Error; err != nil {
		log.Error(c.IP(), " error ", err.Error())
		return nil, err
	}

	for _, rec := range candidates {
		if helper.CompareHashBcript(token, rec.Token) == nil {
			return &rec, nil
		}
	}

	return nil, errors.New("refresh token not found")
}
