package interfaces

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
)

type IImageRepository interface {
	Update(ctx context.Context, dto *dtos.UpdateImageDto) error
	Create(ctx context.Context, dto *dtos.CreateImageDto) (int, error)
	Delete(ctx context.Context, id int) (string, error)
}
