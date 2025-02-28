package dto

type AuthUserRegisterRequest struct {
	FullName string `json:"fullName" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Gender   *bool   `json:"gender" validate:"required"`
}
