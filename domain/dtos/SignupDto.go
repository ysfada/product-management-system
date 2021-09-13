package dtos

type SignupDto struct {
	Username string `json:"username" validate:"required,min=2,max=16"`
	Password string `json:"password" validate:"required,min=8"`
}
