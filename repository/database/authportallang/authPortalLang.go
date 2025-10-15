package authportallang

import (
	"github.com/nurefendi/auth-service-using-golang/config/database"
	"github.com/nurefendi/auth-service-using-golang/dto"
	"github.com/nurefendi/auth-service-using-golang/repository/dao"
	"github.com/nurefendi/auth-service-using-golang/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

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

	deleteRecord := db.Delete(&dao.AuthPortalLang{}, "id = ?", portalID)
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
