package usecase

import (
	"auth-service/dto"
	"auth-service/repository/dao"
	portalRepository "auth-service/repository/database/authportal"
	portalLangRepository "auth-service/repository/database/authportallang"
	"auth-service/tools/helper"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func PortalUseCase() Portal {
	return &portalUseCase{}
}

// Delete implements Portal.
func (p *portalUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	_, errf := portalRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return errf
	}
	return portalRepository.Delete(c, id)
}

// Save implements Portal.
func (p *portalUseCase) Save(c *fiber.Ctx, data *dto.PortalDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var languages []dao.AuthPortalLang
	for _, lang := range data.Languages {
		languages = append(languages, dao.AuthPortalLang{
			Name:        lang.PortalName,
			Lang:        lang.LanguageCode,
			Description: &lang.Description,
			AuditorDAO: dao.AuditorDAO{
				CreatedBy: currentAccess.UserAccess.Email,
			},
		})
	}
	saveError := portalRepository.Save(c, &dao.AuthPortal{
		Order:    data.Order,
		Path:     data.Path,
		Icon:     data.Icon,
		FontIcon: data.FontIcon,
		AuditorDAO: dao.AuditorDAO{
			CreatedBy: currentAccess.UserAccess.Email,
		},
		Lang: languages,
	})
	if saveError != nil {
		return saveError
	}
	return &fiber.Error{}
}

// Update implements Portal.
func (p *portalUseCase) Update(c *fiber.Ctx, data *dto.PortalDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if data.ID == nil {
		log.Error(currentAccess.RequestID, " ID must not nil")
		return fiber.NewError(fiber.StatusBadRequest, "Portal id required")
	}
	datadb, errf := portalRepository.FindById(c, *data.ID)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return errf
	}
	existingLangMap := make(map[string]dao.AuthPortalLang)
	for _, lang := range datadb.Lang {
		existingLangMap[lang.Lang] = lang
	}
	var languages []dao.AuthPortalLang
	for _, lang := range data.Languages {
		if existingLang, found := existingLangMap[lang.LanguageCode]; found {
			existingLang.Name = lang.PortalName
			existingLang.Description = &lang.Description
			existingLang.AuditorDAO.ModifiedBy = &currentAccess.UserAccess.Email
			languages = append(languages, existingLang)
			delete(existingLangMap, lang.LanguageCode)
		} else {
			languages = append(languages, dao.AuthPortalLang{
				Name:        lang.PortalName,
				Lang:        lang.LanguageCode,
				Description: &lang.Description,
				AuditorDAO: dao.AuditorDAO{
					CreatedBy: currentAccess.UserAccess.Email,
				},
			})
		}
	}

	for _, langToDelete := range existingLangMap {
		err := portalLangRepository.Delete(c, langToDelete.ID)
		if err != nil {
			log.Error(currentAccess.RequestID, "Failed to delete lang:", langToDelete.ID, "Error:", err.Message)
			return err
		}
	}
	saveError := portalRepository.Save(c, &dao.AuthPortal{
		Order:    data.Order,
		Path:     data.Path,
		Icon:     data.Icon,
		FontIcon: data.FontIcon,
		AuditorDAO: dao.AuditorDAO{
			ModifiedBy: &currentAccess.UserAccess.Email,
			ID:         *data.ID,
		},
		Lang: languages,
	})
	if saveError != nil {
		log.Error(currentAccess.RequestID, "Failed to save portal:", saveError.Message)
		return saveError
	}

	return &fiber.Error{}
}

// FindById implements Portal.
func (p *portalUseCase) FindById(c *fiber.Ctx, id uuid.UUID) (*dto.PortalDto, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, errf := portalRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil, errf
	}
	var result dto.PortalDto
	err := helper.Map(data, &result)
	if err != nil {
		log.Error(currentAccess.RequestID, err.Error())
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Erorr bind data")
	}
	return &result, &fiber.Error{}
}
// FindAll implements Portal.
func (p *portalUseCase) FindAll(c *fiber.Ctx, r dto.PortalPagination) ([]dto.PortalUserDto, int64, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, total, errf := portalRepository.FindAll(c, r)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil,0, errf
	}
	var result []dto.PortalUserDto
	for _, v := range data {
		var portalName, portalDescription string
		for _, w := range v.Lang {
			if w.Lang == currentAccess.LanguageCode {
				portalName = w.Name
				portalDescription = *w.Description
				break
			}
		}
		result = append(result, dto.PortalUserDto{
			ID: v.ID,
			Order: v.Order,
			Path: v.Path,
			Icon: v.Icon,
			FontIcon: v.FontIcon,
			PortalName: portalName,
			Description: portalDescription,
		})
	}
	return result, total, nil
}