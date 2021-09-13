package dtos

type CreateProductVariantDto struct {
	Name      string  `json:"name" validate:"required,min=2,max=16"`
	ProductId int     `json:"product_id" validate:"required,number"`
	Price     float64 `json:"price" validate:"required,number"`
	Stock     int     `json:"stock" validate:"required,number"`
}
