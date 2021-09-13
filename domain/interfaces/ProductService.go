package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/ysfada/product-management-system/domain/dtos"
)

type IProductService interface {
	Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*dtos.ProductPaginatedDto, error)
	GetByID(ctx context.Context, id int) (res *dtos.ProductDto, err error)
	Update(ctx context.Context, dto *dtos.UpdateProductDto) error
	Create(ctx context.Context, dto *dtos.CreateProductDto) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*dtos.ProductPaginatedDto, error)
	GetImages(ctx context.Context, id int) ([]*dtos.ImageDto, error)
	AddImage(ctx context.Context, id int, fileheader *multipart.FileHeader) error
	RemoveImage(ctx context.Context, id int, imageID int) error
	FetchVariants(ctx context.Context, id int, page int, size int, sortBy string, orderBy string) (*dtos.ProductVariantPaginatedDto, error)
	SearchVariants(ctx context.Context, q string, id int, page int, size int, sortBy string, orderBy string, attrs []*dtos.AttributeSearchQueryDto) (*dtos.ProductVariantPaginatedDto, error)
	GetVariantByID(ctx context.Context, id int, variantID int) (*dtos.ProductVariantDto, error)
	CreateVariant(ctx context.Context, dto *dtos.CreateProductVariantDto) error
	UpdateVariant(ctx context.Context, dto *dtos.UpdateProductVariantDto) error
	DeleteVariant(ctx context.Context, id int, variantID int) error
	GetAttributes(ctx context.Context, id int, variantID int) ([]*dtos.AttributeDto, error)
	AddAttribute(ctx context.Context, dto *dtos.CreateProductVariantAttributeDto) error
	RemoveAttribute(ctx context.Context, id int, variantID int, attributeID int) error
}
