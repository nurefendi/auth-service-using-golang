package dto

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type AuthUserRegisterRequest struct {
	FullName string `json:"fullName" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Gender   int    `json:"gender" validate:"required"`
}

type AuthUserLoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
}

type CurrentUserAccess struct {
	jwt.StandardClaims
	UserID   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName" `
}

type UserLocals struct {
	UserAccess   *CurrentUserAccess
	RequestID    string
	LanguageCode string
	ChannelID    string
}
type AuthUserResponse struct {
	UserID     uuid.UUID   `json:"userId"`
	UserName   string      `json:"userName"`
	Email      string      `json:"email"`
	GroupIDs   []uuid.UUID `json:"groupIds"`
	FullName   string      `json:"fullName"`
	Gender     int         `json:"gender"`
	GenderName string      `json:"genderName"`
	Picture    *string     `json:"picture"`
}
type AuthCheckAccessRequest struct {
	Path   string `json:"path" validate:"required"`
	Mathod string `json:"method" validate:"required"`
}
