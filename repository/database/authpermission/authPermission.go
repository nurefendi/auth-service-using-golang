package authpermission

import (
	enums "auth-service/common/enums/httpmethod"
	"auth-service/config/database"
	"auth-service/repository/dao"
	"auth-service/tools/locals"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func Save(c *fiber.Ctx, data *dao.AuthPermission) error {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals(c)
	if db == nil {
		log.Error(currentAcess.RequestID, " error cannot find db connection")
		return errors.New("database connection is nil")
	}
	saveRecord := db.Save(data)
	if saveRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", saveRecord.Error.Error())
		return saveRecord.Error
	}
	log.Info(currentAcess.RequestID, " Save data to DB, Affected rows :", saveRecord.RowsAffected)
	return nil
}

func FindById(c *fiber.Ctx, id uuid.UUID) (*dao.AuthPermission, error) {
	db := database.GetDBConnection(c)
	currentAcess := locals.GetLocals(c)
	data := dao.AuthPermission{}
	findRecord := db.Where("id = ?", id).
		Preload("Function").
		First(&data)
	if findRecord.Error != nil {
		log.Error(currentAcess.RequestID, " error ", findRecord.Error.Error())
		return nil, findRecord.Error
	}

	return &data, nil
}

func FindByGroupIdAndPathAndMethod(c *fiber.Ctx) (bool, error) {
	db := database.GetDBConnection(c)
	currentAccess := locals.GetLocals(c)

	if db == nil {
		log.Error(currentAccess.RequestID, " database connection is nil")
		return false, errors.New("database connection is nil")
	}

	permissions := map[string]string{
		enums.GET.Name():    "ap.grant_read",
		enums.DELETE.Name(): "ap.grant_delete",
		enums.PUT.Name():    "ap.grant_update",
		enums.PATCH.Name():  "ap.grant_update",
		enums.POST.Name():   "ap.grant_create",
	}

	permissionColumn, ok := permissions[c.Method()]
	if !ok {
		return false, errors.New("unsupported HTTP method")
	}

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM auth_permission ap 
			INNER JOIN auth_function af ON ap.function_id = af.id 
			WHERE af.path = ? 
			AND ap.group_id = ? 
			AND ` + permissionColumn + ` = 1
		) 
	`

	var isExist bool
	err := db.Raw(query, c.Path(), currentAccess.UserAccess.GroupID).Scan(&isExist).Error
	if err != nil {
		log.Error(currentAccess.RequestID, " error ", err.Error())
		return false, err
	}

	return isExist, nil
}