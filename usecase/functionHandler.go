package usecase

import (
	"auth-service/dto"
	"auth-service/repository/dao"
	functionRepository "auth-service/repository/database/authfunction"
	functionLangRepository "auth-service/repository/database/authfunctionlang"
	"auth-service/tools/helper"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func FunctionUSeCase() Function {
	return &functionUseCase{}
}

// Delete implements Function.
func (f *functionUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	_, errf := functionRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return errf
	}
	return functionRepository.Delete(c, id)
}

// FindAll implements Function.
func (f *functionUseCase) FindAll(c *fiber.Ctx, r dto.FunctionPagination) ([]dto.FunctionUserDto, int64, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, total, errf := functionRepository.FindAll(c, r)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil,0, errf
	}
	var result []dto.FunctionUserDto
	for _, v := range data {
		var functionName, funcDesc string
		for _, w := range v.Lang {
			if w.Lang == currentAccess.LanguageCode {
				functionName = w.Name
				funcDesc = *w.Description
				break
			}
		}
		result = append(result, dto.FunctionUserDto{
			ID: v.ID,
			Order: v.Order,
			Path: v.Path,
			Icon: v.Icon,
			FontIcon: v.FontIcon,
			FunctionName: functionName,
			Description: funcDesc,
			PortalID: v.PortalID,
			Method: v.Method,
			Position: v.Position,
			ShortcutKey: v.ShortcutKey,
		})
	}
	return result, total, nil
}

// FindById implements Function.
func (f *functionUseCase) FindById(c *fiber.Ctx, id uuid.UUID) (*dto.FunctionDto, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, errf := functionRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil, errf
	}
	var result dto.FunctionDto
	err := helper.Map(data, &result)
	if err != nil {
		log.Error(currentAccess.RequestID, err.Error())
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Erorr bind data")
	}
	return &result, &fiber.Error{}
}

// Save implements Function.
func (f *functionUseCase) Save(c *fiber.Ctx, data *dto.FunctionDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var languages []dao.AuthFunctionLang
	for _, lang := range data.Languages {
		languages = append(languages, dao.AuthFunctionLang{
			Name:        lang.FunctionName,
			Lang:        lang.LanguageCode,
			Description: &lang.Description,
			AuditorDAO: dao.AuditorDAO{
				CreatedBy: currentAccess.UserAccess.Email,
			},
		})
	}
	saveError := functionRepository.Save(c, &dao.AuthFunction{
		ParentID:    data.ParentID,
		PortalID:    data.PortalID,
		Method:      data.Method,
		IsShow:      data.IsShow,
		Position:    data.Position,
		ShortcutKey: data.ShortcutKey,
		Order:       data.Order,
		Path:        data.Path,
		Icon:        data.Icon,
		FontIcon:    data.FontIcon,
		AuditorDAO: dao.AuditorDAO{
			CreatedBy: currentAccess.UserAccess.Email,
		},
		Lang: languages,
	})
	if saveError != nil {
		log.Error(currentAccess.RequestID, saveError.Error())
		return fiber.NewError(fiber.StatusInternalServerError, " Can't save data to db")
	}
	return nil
}

// Update implements Function.
func (f *functionUseCase) Update(c *fiber.Ctx, data *dto.FunctionDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if data.ID == nil {
		log.Error(currentAccess.RequestID, " ID must not nul")
		return fiber.NewError(fiber.StatusBadRequest, "Portal id required")
	}
	datadb, errf := functionRepository.FindById(c, *data.ID)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return errf
	}
	existingLangMap := make(map[string]dao.AuthFunctionLang)
	for _, lang := range datadb.Lang {
		existingLangMap[lang.Lang] = lang
	}
	var languages []dao.AuthFunctionLang
	for _, lang := range data.Languages {
		if existingLang, found := existingLangMap[lang.LanguageCode]; found {
			existingLang.Name = lang.FunctionName
			existingLang.Description = &lang.Description
			existingLang.AuditorDAO.ModifiedBy = &currentAccess.UserAccess.Email
			languages = append(languages, existingLang)
			delete(existingLangMap, lang.LanguageCode)
		} else {
			languages = append(languages, dao.AuthFunctionLang{
				Name:        lang.FunctionName,
				Lang:        lang.LanguageCode,
				Description: &lang.Description,
				AuditorDAO: dao.AuditorDAO{
					CreatedBy: currentAccess.UserAccess.Email,
				},
			})
		}
	}

	for _, langToDelete := range existingLangMap {
		err := functionLangRepository.Delete(c, langToDelete.ID)
		if err != nil {
			log.Error(currentAccess.RequestID, "Failed to delete lang:", langToDelete.ID, "Error:", err.Message)
			return err
		}
	}
	saveError := functionRepository.Save(c, &dao.AuthFunction{
		ParentID:    data.ParentID,
		PortalID:    data.PortalID,
		Method:      data.Method,
		IsShow:      data.IsShow,
		Position:    data.Position,
		ShortcutKey: data.ShortcutKey,
		Order:       data.Order,
		Path:        data.Path,
		Icon:        data.Icon,
		FontIcon:    data.FontIcon,
		AuditorDAO: dao.AuditorDAO{
			CreatedBy: currentAccess.UserAccess.Email,
		},
		Lang: languages,
	})
	if saveError != nil {
		log.Error(currentAccess.RequestID, "Failed to save function:", saveError.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save data")
	}

	return nil
}
