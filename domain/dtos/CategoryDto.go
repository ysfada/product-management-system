package dtos

type CategoryDto struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,max=32"`
	Description string `json:"description"`
}

type CategoryPaginatedDto struct {
	PaginationDto
	Categories []*CategoryDto `json:"categories"`
}

type CategoryProductsPaginatedDto struct {
	PaginationDto
	CategoryDto
	Products []*ProductDto `json:"products"`
}
