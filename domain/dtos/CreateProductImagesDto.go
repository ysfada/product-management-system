package dtos

type CreateProductImagesDto struct {
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}
