package entities

type Product struct {
	ID              int               `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	CategoryID      int               `json:"category_id"`
	Category        *Category         `json:"category"`
	Images          []*Image          `json:"images"`
	ProductVariants []*ProductVariant `json:"product_variants"`
	Timestamps
}

type ProductPaginated struct {
	Pagination
	Products []*Product `json:"products"`
}
