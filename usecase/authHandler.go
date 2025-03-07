package usecase

import (
	"auth-service/dto"
	jwtMidleware "auth-service/middleware/jwt"
	"auth-service/tools/helper"
	"auth-service/tools/locals"
	"strings"
	"time"

	"auth-service/repository/dao"
	authUserRepository "auth-service/repository/database/authuser"
	authUserGroupRepository "auth-service/repository/database/authusergroup"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func AuthUSeCase() Auth {
	return &authUseCase{}
}

// Login implements Auth.
func (a *authUseCase) Login(c *fiber.Ctx) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	payload := locals.GetLocals[dto.AuthUserLoginRequest](c,locals.PayloadLocalKey)
	log.Info(currentAccess.RequestID, " do login ", payload.Email)
	data, errr := authUserRepository.FindByEmail(c, &payload.Email)
	if errr != nil {
		log.Error(currentAccess.RequestID, errr.Message)
		return errr
	}
	if data.ID == uuid.Nil {
		log.Error(currentAccess.RequestID, " Data not found!")
		return fiber.NewError(fiber.StatusBadRequest, "Wrong Email!")
	}
	err := helper.CompareHashBcript(payload.Password, data.Password)
	if err != nil {
		log.Error(currentAccess.RequestID, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Wrong Password!")
	}
	userGroups, errr := authUserGroupRepository.FindByUserId(c, data.ID)
	if errr != nil {
		log.Error(currentAccess.RequestID, errr.Message)
		return errr
	}

	groupIds := make([]uuid.UUID, len(*userGroups))
	for _, v := range *userGroups {
		groupIds = append(groupIds, v.GroupID)
	}
	expirationTime := time.Now().Add(15 * time.Minute)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

	payloadGenerateToken := dto.CurrentUserAccess{
		UserID: data.ID,
		UserName: data.Username,
		Email: data.Email,
		GroupIDs: groupIds,
		FullName: data.FullName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token, err := jwtMidleware.GenerateToken(payloadGenerateToken)

	if err != nil {
		log.Error(currentAccess.RequestID, err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed generate token!")
	}

	payloadGenerateToken.StandardClaims = jwt.StandardClaims{
		ExpiresAt: refreshExpiration.Unix(),
	}
	refreshToken, err := jwtMidleware.GenerateToken(payloadGenerateToken)
	if err != nil {
		log.Error(currentAccess.RequestID, err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed generate refresh token!")
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    *token,
		Expires:  expirationTime,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshToken,
		Expires:  refreshExpiration,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})
	return nil
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

func (a *authUseCase) CheckUserNameExist(c *fiber.Ctx, userNames *string) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, err := authUserRepository.FindByUserName(c, userNames)
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