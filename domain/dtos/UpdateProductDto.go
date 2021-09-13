package dtos

type UpdateProductDto struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,max=16"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id" validate:"required,number"`
}
