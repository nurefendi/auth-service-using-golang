package authportal

import (
	"auth-service/config/database"
	"auth-service/dto"
	"auth-service/repository/dao"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func Save(c *fiber.Ctx, data *dao.AuthPortal) *fiber.Error {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Unable to conect database",
		}
	}
	saveRecord := db.Create(&data)
	if saveRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", saveRecord.Error.Error())
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: saveRecord.Error.Error(),
		}
	}
	log.Info(currentAcess.RequestID, " Save data to DB, Affected rows :", saveRecord.RowsAffected)
	return nil
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

	deleteRecord := db.Delete(&dao.AuthPortal{}, "id = ?", portalID)
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
func FindById(c *fiber.Ctx, portalID uuid.UUID) (dao.AuthPortal, *fiber.Error) {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var result dao.AuthPortal
	if db == nil {
		log.Error(currentAccess.RequestID, " error cannot find db connection")
		return dao.AuthPortal{}, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Unable to connect to database",
		}
	}

	if err := db.Where("id = ?", portalID).
		Preload("AuthPortal").
		First(&result).Error; err != nil {
		log.Error(currentAccess.RequestID, " Data not found")
		return dao.AuthPortal{}, &fiber.Error{
			Code:    fiber.StatusNoContent,
			Message: "Data not found",
		}
	}

	return result, &fiber.Error{}

}
func FindAll(c *fiber.Ctx, r dto.PortalPagination) ([]dao.AuthPortal, int64, *fiber.Error)  {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var authPortals []dao.AuthPortal
	var total int64
	query := db.Model(&dao.AuthPortal{})

	if r.Search != "" {
		searchPattern := "%" + r.Search + "%"
		query = query.Where("path LIKE ? OR icon LIKE ?", searchPattern, searchPattern)
	}
	
	query.Count(&total)
	offset := (r.Offset - 1) * r.Limit

	result := query.Preload("Lang").Order("order ASC").Limit(r.Limit).Offset(offset).Find(&authPortals)
	if result.Error != nil {
		log.Error(currentAccess.RequestID, " error ", result.Error.Error())
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return  authPortals, total, nil
}
