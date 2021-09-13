package interfaces

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/entities"
)

type IAttributeRepository interface {
	Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*entities.AttributePaginated, error)
	GetByID(ctx context.Context, id int) (res *entities.Attribute, err error)
	Update(ctx context.Context, dto *dtos.UpdateAttributeDto) error
	Create(ctx context.Context, dto *dtos.CreateAttributeDto) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*entities.AttributePaginated, error)
}
