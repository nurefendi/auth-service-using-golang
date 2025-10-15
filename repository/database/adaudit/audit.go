package adaudit

import (
	"auth-service/config/database"
	"auth-service/dto"
	"auth-service/repository/dao"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func Save(c *fiber.Ctx, action string, userId *uuid.UUID, ip *string, userAgent *string, metadata *string) error {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return c.Status(fiber.StatusInternalServerError).
			SendString("Database error")
	}

	audit := dao.AuthAudit{
		Action: action,
		UserID: userId,
		IP: ip,
		UserAgent: userAgent,
		Metadata: metadata,
		AuditorDAO: dao.AuditorDAO{
			ID: uuid.New(),
			CreatedBy: currentAcess.RequestID,
		},
	}

	saveRecord := db.Create(&audit)
	if saveRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", saveRecord.Error.Error())
		return saveRecord.Error
	}
	log.Info(currentAcess.RequestID, " Save audit to DB, Affected rows :", saveRecord.RowsAffected)
	return nil
}
