package dto

import (
	"github.com/google/uuid"
)

type UserDto struct {
	ID        *uuid.UUID  `json:"id"`
	GroupIDs  []uuid.UUID `json:"groupIds" validate:"required"`
	Gender    int         `json:"gender" validate:"required"`
	FullName  string      `json:"fullName" validate:"required,min=2"`
	Email     string      `json:"email" validate:"required,email"`
	Username  string      `json:"userName" validate:"required,min=3"`
	Password  string      `json:"password"  validate:"required,min=8"`
	Telephone *string     `json:"telephone" validate:"omitempty"`
	Picture   *string     `json:"picture" validate:"omitempty"`
}

type UserPagination struct {
	PaginationParam
}
