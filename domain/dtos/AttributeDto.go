package dtos

type AttributeDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type AttributePaginatedDto struct {
	PaginationDto
	Attributes []*AttributeDto `json:"attributes"`
}
