package usecase

import (
	"auth-service/dto"
	"auth-service/tools/helper"
	"auth-service/tools/locals"
	"strings"

	"auth-service/repository/dao"
	authUserRepository "auth-service/repository/database/authuser"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func AuthUSeCase() Auth {
	return &authUseCase{}
}

// Login implements Auth.
func (a *authUseCase) Login(c *fiber.Ctx) *fiber.Error {
	panic("unimplemented")
}

// Register implements Auth.
func (a *authUseCase) Register(c *fiber.Ctx) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	payload := locals.GetLocals[dto.AuthUserRegisterRequest](c,locals.PayloadLocalKey)
	log.Info(currentAccess.RequestID, " do register ", payload.Email)
	userNames := strings.Split(payload.Email, "@")
	password, err := helper.HashBcript(payload.Password)
	if err != nil {
		log.Error(currentAccess.RequestID, " Invalid hash password")
		return &fiber.Error{
			Code: fiber.StatusInternalServerError,
			Message: "Invalid hash password",
		}
	}
	fibererr := a.CheckEmailExist(c, &payload.Email)
	if fibererr != nil {
		return &fiber.Error{
			Code: fibererr.Code,
			Message: fibererr.Error(),
		}
	}
	c.Locals(locals.Entity, dao.AuthUser{
		FullName: payload.FullName,
		Email: payload.Email,
		Password: password,
		Gender: payload.Gender,
		HasDeleted: false,
		Username: userNames[0] + userNames[1],
		AuditorDAO: dao.AuditorDAO{
			ID: uuid.New(),
			CreatedBy: c.IP(),
		},
	})
	fibererr = authUserRepository.Save(c)
	if fibererr != nil {
		return fibererr
	}
	return nil
}

// CheckEmailExist implements Auth.
func (a *authUseCase) CheckEmailExist(c *fiber.Ctx, email *string) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, err := authUserRepository.FindByEmail(c, email)
	if err != nil {
		log.Error(currentAccess.RequestID, err)
		return &fiber.Error{
			Code: fiber.StatusInternalServerError,
			Message: err.Error(),
		}
		
	}
	if data.ID != uuid.Nil {
		log.Error(currentAccess.RequestID, " email exist ")
		return &fiber.Error{
			Code: fiber.StatusUnprocessableEntity,
			Message: "email exist",
		}
	}
	return nil
}