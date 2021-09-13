package dtos

type UpdateUserDto struct {
	ID       int    `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=2,max=16"`
	Password string `json:"password" validate:"required,min=8"`
}
