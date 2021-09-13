package dtos

type ProductDto struct {
	ID          int                  `json:"id" validate:"required"`
	Name        string               `json:"name" validate:"required,min=2,max=16"`
	Description string               `json:"description"`
	CategoryID  int                  `json:"category_id" validate:"required,number"`
	Category    *CategoryDto         `json:"category,omitempty"`
	Images      []*ImageDto          `json:"images,omitempty"`
	Variants    []*ProductVariantDto `json:"variants,omitempty"`
}

type ProductPaginatedDto struct {
	PaginationDto
	Products []*ProductDto `json:"products"`
}
