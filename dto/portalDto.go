package dto

type PortalSaveRequest struct {
	Path      string                  `json:"path" validate:"required"`
	Icon      *string                 `json:"icon" validate:"omitempty"`
	FontIcon  *string                 `json:"fontIcon" validate:"omitempty"`
	Languages []PortalLangSaveRequest `json:"languages" validate:"required,dive"`
}

type PortalLangSaveRequest struct {
	PortalName   string `json:"portalName" validate:"required,max=255"`
	Description  string `json:"description" validate:"required"`
	LanguageCode string `json:"languageCode" validate:"required,min=2,max=4"`
}
