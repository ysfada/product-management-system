package dtos

type ChangeUsernameDto struct {
	Username string `json:"username" validate:"required,min=2,max=16"`
}
