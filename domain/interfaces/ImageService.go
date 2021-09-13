package interfaces

import (
	"context"
	"mime/multipart"
)

type IImageService interface {
	Save(ctx context.Context, id int, fileheader *multipart.FileHeader) (int, error)
	Remove(ctx context.Context, imageID int) error
}
