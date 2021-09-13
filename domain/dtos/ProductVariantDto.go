package dtos

type ProductVariantDto struct {
	ID         int             `json:"id" validate:"required"`
	Name       string          `json:"name" validate:"required,min=2,max=16"`
	ProductId  int             `json:"product_id" validate:"required,number"`
	Product    *ProductDto     `json:"product,omitempty"`
	Price      float64         `json:"price" validate:"required,number"`
	Stock      int             `json:"stock" validate:"required,number"`
	Attributes []*AttributeDto `json:"attributes"`
}

type ProductVariantPaginatedDto struct {
	PaginationDto
	ProductDto
	ProductVariants []*ProductVariantDto `json:"product_variants"`
}
