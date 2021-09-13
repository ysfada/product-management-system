package dtos

type CreateProductDto struct {
	Name        string `json:"name" validate:"required,min=2,max=16"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id" validate:"required,number"`
}
