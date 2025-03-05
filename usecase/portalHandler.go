package usecase

import (
	"auth-service/dto"
	"auth-service/repository/dao"
	portalRepository "auth-service/repository/database/authportal"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PortalUseCase() Portal {
	return &portalUseCase{}
}

// Delete implements Portal.
func (p *portalUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	return portalRepository.Delete(c, id)
}

// Save implements Portal.
func (p *portalUseCase) Save(c *fiber.Ctx, data *dto.PortalSaveRequest) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	var languages []dao.AuthPortalLang
	for _, lang := range data.Languages {
		languages = append(languages, dao.AuthPortalLang{
			Name: lang.PortalName,
			Lang: lang.LanguageCode,
			Description: &lang.Description,
			AuditorDAO: dao.AuditorDAO{
				CreatedBy: currentAccess.UserAccess.Email,
			},
		})
	}
	saveError := portalRepository.Save(c, &dao.AuthPortal{
		Order: data.Order,
		Path: data.Path,
		Icon: data.Icon,
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
func (p *portalUseCase) Update(c *fiber.Ctx, data *dto.PortalSaveRequest) *fiber.Error {
	panic("unimplemented")
}
