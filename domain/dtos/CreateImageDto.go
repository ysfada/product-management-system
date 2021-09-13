package dtos

type CreateImageDto struct {
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}
