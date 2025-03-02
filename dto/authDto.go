package dto

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type AuthUserRegisterRequest struct {
	FullName string `json:"fullName" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Gender   *bool  `json:"gender" validate:"required"`
}

type CurrentUserAccess struct {
	jwt.StandardClaims
	UserID   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
	Email    string    `json:"email"`
	GroupID  uuid.UUID `json:"groupId"`
}

type UserLocals struct {
	UserAccess   CurrentUserAccess
	RequestID    string
	LanguageCode string
	ChannelID    string
}
