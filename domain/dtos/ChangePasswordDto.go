package dtos

type ChangePasswordDto struct {
	Password string `json:"password" validate:"required,min=8"`
}
