package dto

import (
	"github.com/google/uuid"
)

type GroupDto struct {
	ID          *uuid.UUID `json:"id"`
	Name        string     `json:"name" validate:"required,min=2"`
	Description *string    `json:"description" validate:"omitempty"`
}

type GroupPagination struct {
	PaginationParam
}
