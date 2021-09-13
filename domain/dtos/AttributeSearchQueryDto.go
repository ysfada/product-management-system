package dtos

type AttributeSearchQueryDto struct {
	Names []string `json:"names"`
	Type  string   `json:"type"`
}
