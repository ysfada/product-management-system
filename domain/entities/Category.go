package entities

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Timestamps
}

type CategoryPaginated struct {
	Pagination
	Categories []*Category `json:"categories"`
}

type CategoryProductsPaginated struct {
	Pagination
	Category
	Products []*Product `json:"products"`
}
