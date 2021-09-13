package entities

type Attribute struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Timestamps
}

type AttributePaginated struct {
	Pagination
	Attributes []*Attribute `json:"attributes"`
}
