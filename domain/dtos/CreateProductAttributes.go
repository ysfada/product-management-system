package dtos

type CreateProductAttributesDto struct {
	ProductID   int `json:"product_id" validate:"required"`
	AttributeID int `json:"attribute_id" validate:"required"`
}
