package interfaces

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
)

type IAttributeService interface {
	Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*dtos.AttributePaginatedDto, error)
	GetByID(ctx context.Context, id int) (res *dtos.AttributeDto, err error)
	Update(ctx context.Context, dto *dtos.UpdateAttributeDto) error
	Create(ctx context.Context, dto *dtos.CreateAttributeDto) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*dtos.AttributePaginatedDto, error)
}
