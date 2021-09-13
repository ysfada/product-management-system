package dtos

type SigninDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
