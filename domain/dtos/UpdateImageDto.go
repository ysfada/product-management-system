package dtos

type UpdateImageDto struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}
