package interfaces

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/entities"
)

type IProductRepository interface {
	Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*entities.ProductPaginated, error)
	GetByID(ctx context.Context, id int) (*entities.Product, error)
	Update(ctx context.Context, dto *dtos.UpdateProductDto) error
	Create(ctx context.Context, dto *dtos.CreateProductDto) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*entities.ProductPaginated, error)
	GetImages(ctx context.Context, id int) ([]*entities.Image, error)
	AddImage(ctx context.Context, id int, imageID int) error
	RemoveImage(ctx context.Context, id int, imageID int) error
	FetchVariants(ctx context.Context, id int, page int, size int, sortBy string, orderBy string) (*entities.ProductVariantPaginated, error)
	SearchVariants(ctx context.Context, q string, id int, page int, size int, sortBy string, orderBy string, attrs []*dtos.AttributeSearchQueryDto) (*entities.ProductVariantPaginated, error)
	GetVariantByID(ctx context.Context, id int, variantID int) (*entities.ProductVariant, error)
	CreateVariant(ctx context.Context, dto *dtos.CreateProductVariantDto) error
	UpdateVariant(ctx context.Context, dto *dtos.UpdateProductVariantDto) error
	DeleteVariant(ctx context.Context, id int, variantID int) error
	GetAttributes(ctx context.Context, id int, variantID int) ([]*entities.Attribute, error)
	AddAttribute(ctx context.Context, dto *dtos.CreateProductVariantAttributeDto) error
	RemoveAttribute(ctx context.Context, id int, variantID int, attributeID int) error
}
