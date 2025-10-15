package authgroup

import (
	"github.com/nurefendi/auth-service-using-golang/config/database"
	"github.com/nurefendi/auth-service-using-golang/dto"
	"github.com/nurefendi/auth-service-using-golang/repository/dao"
	"github.com/nurefendi/auth-service-using-golang/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func Save(c *fiber.Ctx, data *dao.AuthGroup) error {
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

func FindById(c *fiber.Ctx, id uuid.UUID) (*dao.AuthGroup, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data := dao.AuthGroup{}
	findRecord := db.Where("id = ?", id).
		Preload("Lang").
		First(&data)
	if findRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", findRecord.Error.Error())
		return nil, fiber.NewError(fiber.StatusInternalServerError, findRecord.Error.Error())
	}

	return &data, nil
}
func Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)

	if db == nil {
		log.Error(currentAccess.RequestID, " error cannot find db connection")
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Unable to connect to database",
		}
	}

	deleteRecord := db.Delete(&dao.AuthGroup{}, "id = ?", id)
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
func FindAll(c *fiber.Ctx, r dto.GroupPagination) ([]dao.AuthGroup, int64, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var resultDao []dao.AuthGroup
	var total int64
	query := db.Model(&dao.AuthGroup{})

	if r.Search != "" {
		searchPattern := "%" + r.Search + "%"
		query = query.Where("path LIKE ? OR icon LIKE ?", searchPattern, searchPattern)
	}

	query.Count(&total)
	offset := (r.Offset - 1) * r.Limit

	result := query.Preload("Lang").Order("order ASC").Limit(r.Limit).Offset(offset).Find(&resultDao)
	if result.Error != nil {
		log.Error(currentAccess.RequestID, " error ", result.Error.Error())
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return resultDao, total, nil
}
