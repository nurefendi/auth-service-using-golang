package authusergroup

import (
	"auth-service/config/database"
	"auth-service/dto"
	"auth-service/repository/dao"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func FindByUserId(c *fiber.Ctx, userId uuid.UUID) (*[]dao.AuthUserGroup, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var result []dao.AuthUserGroup
	err := db.Model(dao.AuthUserGroup{}).
		Where("user_id = ?", userId).
		Find(&result).Error
	if err != nil {
		log.Info(currentAcess.RequestID, err.Error())
		return nil, fiber.NewError(fiber.StatusBadRequest, "Group not found")
	}
	return &result, nil
}


func Save(c *fiber.Ctx, data dao.AuthUserGroup) *fiber.Error {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
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