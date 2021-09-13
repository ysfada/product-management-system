package interfaces

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
)

type ICategoryService interface {
	Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*dtos.CategoryPaginatedDto, error)
	GetByID(ctx context.Context, id int) (res *dtos.CategoryDto, err error)
	Update(ctx context.Context, dto *dtos.UpdateCategoryDto) error
	Create(ctx context.Context, dto *dtos.CreateCategoryDto) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*dtos.CategoryPaginatedDto, error)
	GetProducts(ctx context.Context, id int, page int, size int, sortBy string, orderBy string) (*dtos.CategoryProductsPaginatedDto, error)
}
