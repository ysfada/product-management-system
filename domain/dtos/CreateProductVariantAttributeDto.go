package dtos

type CreateProductVariantAttributeDto struct {
	ProductVariantID int `json:"product_variant_id"`
	AttributeID      int `json:"attribute_id"`
}
