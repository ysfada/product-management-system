package dtos

type UpdateAttributeDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
