package dto

import "github.com/google/uuid"

type PortalDto struct {
	ID        *uuid.UUID      `json:"id"`
	Order     int             `json:"order" validate:"required"`
	Path      string          `json:"path" validate:"required"`
	Icon      *string         `json:"icon" validate:"omitempty"`
	FontIcon  *string         `json:"fontIcon" validate:"omitempty"`
	Languages []PortalLangDto `json:"languages" validate:"required,dive"`
}

type PortalLangDto struct {
	ID           *uuid.UUID `json:"id"`
	PortalName   string     `json:"portalName" validate:"required,max=255"`
	Description  string     `json:"description" validate:"required"`
	LanguageCode string     `json:"languageCode" validate:"required,min=2,max=4"`
}
type PortalUserDto struct {
	ID          uuid.UUID `json:"id"`
	Order       int       `json:"order" validate:"required"`
	Path        string    `json:"path" validate:"required"`
	Icon        *string   `json:"icon" validate:"omitempty"`
	FontIcon    *string   `json:"fontIcon" validate:"omitempty"`
	PortalName  string    `json:"portalName" validate:"required,max=255"`
	Description string    `json:"description" validate:"required"`
}
