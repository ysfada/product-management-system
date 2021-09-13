package services

import (
	"context"
	"mime/multipart"

	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type ProductService struct {
	repository   interfaces.IProductRepository
	imageService interfaces.IImageService
}

var _ interfaces.IProductService = (*ProductService)(nil)

func NewProductService(repository interfaces.IProductRepository, imageService interfaces.IImageService) *ProductService {
	return &ProductService{
		repository:   repository,
		imageService: imageService,
	}
}

func (s *ProductService) Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*dtos.ProductPaginatedDto, error) {
	if products, err := s.repository.Fetch(ctx, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var productsDto dtos.ProductPaginatedDto
		for _, product := range products.Products {
			productDto := &dtos.ProductDto{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				CategoryID:  product.CategoryID,
				Category: &dtos.CategoryDto{
					ID:          product.Category.ID,
					Name:        product.Category.Name,
					Description: product.Category.Description,
				},
				// Variants:    []*dtos.ProductVariantDto{},
			}

			for _, image := range product.Images {
				imageDto := &dtos.ImageDto{
					ID:           image.ID,
					Name:         image.Name,
					ImageUrl:     image.ImageUrl,
					ThumbnailUrl: image.ThumbnailUrl,
				}

				productDto.Images = append(productDto.Images, imageDto)
			}

			productsDto.Products = append(productsDto.Products, productDto)
		}

		productsDto.TotalPage = products.TotalPage
		productsDto.CurrentPage = products.CurrentPage
		productsDto.NextPage = products.NextPage
		productsDto.PreviousPage = products.PreviousPage
		productsDto.Count = products.Count
		productsDto.Size = products.Size

		return &productsDto, nil
	}
}

func (s *ProductService) GetByID(ctx context.Context, id int) (res *dtos.ProductDto, err error) {
	if product, err := s.repository.GetByID(ctx, id); err != nil {
		return nil, err
	} else {
		if product != nil {
			productDto := &dtos.ProductDto{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				CategoryID:  product.CategoryID,
				Category: &dtos.CategoryDto{
					ID:          product.Category.ID,
					Name:        product.Category.Name,
					Description: product.Category.Description,
				},
				// Variants:    []*dtos.ProductVariantDto{},
			}

			for _, image := range product.Images {
				imageDto := &dtos.ImageDto{
					ID:           image.ID,
					Name:         image.Name,
					ImageUrl:     image.ImageUrl,
					ThumbnailUrl: image.ThumbnailUrl,
				}

				productDto.Images = append(productDto.Images, imageDto)
			}

			return productDto, nil
		}
		return nil, nil
	}
}

func (s *ProductService) Update(ctx context.Context, dto *dtos.UpdateProductDto) error {
	return s.repository.Update(ctx, dto)
}

func (s *ProductService) Create(ctx context.Context, dto *dtos.CreateProductDto) error {
	return s.repository.Create(ctx, dto)
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func (s *ProductService) Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*dtos.ProductPaginatedDto, error) {
	if products, err := s.repository.Search(ctx, q, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var productsDto dtos.ProductPaginatedDto
		for _, product := range products.Products {
			productDto := &dtos.ProductDto{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				CategoryID:  product.CategoryID,
				Category: &dtos.CategoryDto{
					ID:          product.Category.ID,
					Name:        product.Category.Name,
					Description: product.Category.Description,
				},
				// Variants:    []*dtos.ProductVariantDto{},
			}

			for _, image := range product.Images {
				imageDto := &dtos.ImageDto{
					ID:           image.ID,
					Name:         image.Name,
					ImageUrl:     image.ImageUrl,
					ThumbnailUrl: image.ThumbnailUrl,
				}

				productDto.Images = append(productDto.Images, imageDto)
			}

			productsDto.Products = append(productsDto.Products, productDto)
		}

		productsDto.TotalPage = products.TotalPage
		productsDto.CurrentPage = products.CurrentPage
		productsDto.NextPage = products.NextPage
		productsDto.PreviousPage = products.PreviousPage
		productsDto.Count = products.Count
		productsDto.Size = products.Size

		return &productsDto, nil
	}
}

func (s *ProductService) GetImages(ctx context.Context, id int) ([]*dtos.ImageDto, error) {
	if images, err := s.repository.GetImages(ctx, id); err != nil {
		return nil, err
	} else {
		var imagesDto []*dtos.ImageDto
		for _, image := range images {
			imageDto := &dtos.ImageDto{
				ID:           image.ID,
				Name:         image.Name,
				ImageUrl:     image.ImageUrl,
				ThumbnailUrl: image.ThumbnailUrl,
			}

			imagesDto = append(imagesDto, imageDto)
		}

		return imagesDto, nil
	}
}

func (s *ProductService) AddImage(ctx context.Context, id int, fileheader *multipart.FileHeader) error {
	if imageID, err := s.imageService.Save(ctx, id, fileheader); err != nil {
		return err
	} else {
		return s.repository.AddImage(ctx, id, imageID)
	}
}

func (s *ProductService) RemoveImage(ctx context.Context, id int, imageID int) error {
	if err := s.repository.RemoveImage(ctx, id, imageID); err != nil {
		return err
	} else {
		return s.imageService.Remove(ctx, imageID)
	}
}

func (s *ProductService) FetchVariants(ctx context.Context, id int, page int, size int, sortBy string, orderBy string) (*dtos.ProductVariantPaginatedDto, error) {
	if productVariants, err := s.repository.FetchVariants(ctx, id, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var productVariantsDto dtos.ProductVariantPaginatedDto
		for _, variant := range productVariants.ProductVariants {
			productVariantDto := &dtos.ProductVariantDto{
				ID:        variant.ID,
				Name:      variant.Name,
				ProductId: variant.ProductId,
				// Product:    &dtos.ProductDto{},
				Price: variant.Price,
				Stock: variant.Stock,
			}

			for _, attribute := range variant.Attributes {
				attributeDto := &dtos.AttributeDto{
					ID:   attribute.ID,
					Name: attribute.Name,
					Type: attribute.Type,
				}
				productVariantDto.Attributes = append(productVariantDto.Attributes, attributeDto)
			}

			productVariantsDto.ProductVariants = append(productVariantsDto.ProductVariants, productVariantDto)
		}

		productVariantsDto.TotalPage = productVariants.TotalPage
		productVariantsDto.CurrentPage = productVariants.CurrentPage
		productVariantsDto.NextPage = productVariants.NextPage
		productVariantsDto.PreviousPage = productVariants.PreviousPage
		productVariantsDto.Count = productVariants.Count
		productVariantsDto.Size = productVariants.Size

		return &productVariantsDto, nil
	}
}

func (s *ProductService) GetVariantByID(ctx context.Context, id int, variantID int) (*dtos.ProductVariantDto, error) {
	if productVariant, err := s.repository.GetVariantByID(ctx, id, variantID); err != nil {
		return nil, err
	} else {
		if productVariant != nil {
			productVariantDto := &dtos.ProductVariantDto{
				ID:        productVariant.ID,
				Name:      productVariant.Name,
				ProductId: productVariant.ProductId,
				Product: &dtos.ProductDto{
					ID:          productVariant.Product.ID,
					Name:        productVariant.Product.Name,
					Description: productVariant.Product.Description,
					CategoryID:  productVariant.Product.CategoryID,
					Category: &dtos.CategoryDto{
						ID:          productVariant.Product.Category.ID,
						Name:        productVariant.Product.Category.Name,
						Description: productVariant.Product.Category.Description,
					},
					// Images:      []*dtos.ImageDto{},
					// Variants:    []*dtos.ProductVariantDto{},
				},
				Price: productVariant.Price,
				Stock: productVariant.Stock,
			}

			for _, attribute := range productVariant.Attributes {
				attributeDto := &dtos.AttributeDto{
					ID:   attribute.ID,
					Name: attribute.Name,
					Type: attribute.Type,
				}

				productVariantDto.Attributes = append(productVariantDto.Attributes, attributeDto)
			}

			return productVariantDto, nil
		}
		return nil, nil
	}
}

func (s *ProductService) CreateVariant(ctx context.Context, dto *dtos.CreateProductVariantDto) error {
	return s.repository.CreateVariant(ctx, dto)
}

func (s *ProductService) UpdateVariant(ctx context.Context, dto *dtos.UpdateProductVariantDto) error {
	return s.repository.UpdateVariant(ctx, dto)
}

func (s *ProductService) DeleteVariant(ctx context.Context, id int, variantID int) error {
	return s.repository.DeleteVariant(ctx, id, variantID)
}

func (s *ProductService) GetAttributes(ctx context.Context, id int, variantID int) ([]*dtos.AttributeDto, error) {
	if attributes, err := s.repository.GetAttributes(ctx, id, variantID); err != nil {
		return nil, err
	} else {
		var attributesDto []*dtos.AttributeDto
		for _, attribute := range attributes {
			attributeDto := &dtos.AttributeDto{
				ID:   attribute.ID,
				Name: attribute.Name,
				Type: attribute.Type,
			}

			attributesDto = append(attributesDto, attributeDto)
		}

		return attributesDto, nil
	}
}

func (s *ProductService) AddAttribute(ctx context.Context, dto *dtos.CreateProductVariantAttributeDto) error {
	return s.repository.AddAttribute(ctx, dto)
}

func (s *ProductService) RemoveAttribute(ctx context.Context, id int, variantID int, attributeID int) error {
	return s.repository.RemoveAttribute(ctx, id, variantID, attributeID)
}

func (s *ProductService) SearchVariants(ctx context.Context, q string, id int, page int, size int, sortBy string, orderBy string, attrs []*dtos.AttributeSearchQueryDto) (*dtos.ProductVariantPaginatedDto, error) {
	if productVariants, err := s.repository.SearchVariants(ctx, q, id, page, size, sortBy, orderBy, attrs); err != nil {
		return nil, err
	} else {
		var productVariantsDto dtos.ProductVariantPaginatedDto
		for _, variant := range productVariants.ProductVariants {
			productVariantDto := &dtos.ProductVariantDto{
				ID:        variant.ID,
				Name:      variant.Name,
				ProductId: variant.ProductId,
				// Product:    &dtos.ProductDto{},
				Price: variant.Price,
				Stock: variant.Stock,
			}

			for _, attribute := range variant.Attributes {
				attributeDto := &dtos.AttributeDto{
					ID:   attribute.ID,
					Name: attribute.Name,
					Type: attribute.Type,
				}
				productVariantDto.Attributes = append(productVariantDto.Attributes, attributeDto)
			}

			productVariantsDto.ProductVariants = append(productVariantsDto.ProductVariants, productVariantDto)
		}

		productVariantsDto.ID = productVariants.ID
		productVariantsDto.Name = productVariants.Name
		productVariantsDto.Description = productVariants.Description
		productVariantsDto.CategoryID = productVariants.CategoryID
		productVariantsDto.Category = &dtos.CategoryDto{
			ID:          productVariants.Category.ID,
			Name:        productVariants.Category.Name,
			Description: productVariants.Category.Description,
		}

		productVariantsDto.TotalPage = productVariants.TotalPage
		productVariantsDto.CurrentPage = productVariants.CurrentPage
		productVariantsDto.NextPage = productVariants.NextPage
		productVariantsDto.PreviousPage = productVariants.PreviousPage
		productVariantsDto.Count = productVariants.Count
		productVariantsDto.Size = productVariants.Size

		return &productVariantsDto, nil
	}
}
