package dtos

type UpdateCategoryDto struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,max=32"`
	Description string `json:"description"`
}
