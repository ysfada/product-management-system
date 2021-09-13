package entities

type ProductVariant struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	ProductId  int          `json:"product_id"`
	Product    *Product     `json:"product"`
	Price      float64      `json:"price"`
	Stock      int          `json:"stock"`
	Attributes []*Attribute `json:"attributes"`
	Timestamps
}

type ProductVariantPaginated struct {
	Pagination
	Product
	ProductVariants []*ProductVariant `json:"product_variants"`
}
