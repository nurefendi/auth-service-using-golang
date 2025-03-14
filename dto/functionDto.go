package dto

import (
	enums "auth-service/common/enums/httpmethod"

	"github.com/google/uuid"
)

type FunctionDto struct {
	ID          *uuid.UUID        `json:"id"`
	PortalID    uuid.UUID         `json:"portalId" validate:"required"`
	Method      enums.HttpMethod  `json:"method" validate:"required"`
	ParentID    *uuid.UUID        `json:"parentId" validate:"omitempty"`
	Order       int               `json:"order" validate:"required"`
	Position    string            `json:"position" validate:"required"`
	Path        string            `json:"path" validate:"required"`
	Icon        *string           `json:"icon" validate:"omitempty"`
	FontIcon    *string           `json:"fontIcon" validate:"omitempty"`
	ShortcutKey *string           `json:"shortcutKey" validate:"omitempty"`
	IsShow      bool              `json:"isShow" validate:"required"`
	Languages   []FunctionLangDto `json:"languages" validate:"required,dive"`
}

type FunctionLangDto struct {
	ID           *uuid.UUID `json:"id"`
	FunctionName string     `json:"functionName" validate:"required,max=255"`
	Description  string     `json:"description" validate:"required"`
	LanguageCode string     `json:"languageCode" validate:"required,min=2,max=4"`
}
type FunctionUserDto struct {
	ID           uuid.UUID        `json:"id"`
	PortalID     uuid.UUID        `json:"portalId" validate:"required"`
	Method       enums.HttpMethod `json:"method" validate:"required"`
	Order        int              `json:"order" validate:"required"`
	Position     string           `json:"position" validate:"required"`
	Path         string           `json:"path" validate:"required"`
	Icon         *string          `json:"icon" validate:"omitempty"`
	FontIcon     *string          `json:"fontIcon" validate:"omitempty"`
	ShortcutKey  *string          `json:"shortcutKey" validate:"omitempty"`
	IsShow       bool             `json:"isShow" validate:"required"`
	FunctionName string           `json:"functionName" validate:"required,max=255"`
	Description  string           `json:"description" validate:"required"`
}
type FunctionPagination struct {
	PaginationParam
}
