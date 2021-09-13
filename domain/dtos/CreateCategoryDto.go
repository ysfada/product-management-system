package dtos

type CreateCategoryDto struct {
	Name        string `json:"name" validate:"required,min=2,max=32"`
	Description string `json:"description"`
}
