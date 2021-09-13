package services

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type CategoryService struct {
	repository interfaces.ICategoryRepository
}

var _ interfaces.ICategoryService = (*CategoryService)(nil)

func NewCategoryService(repository interfaces.ICategoryRepository) *CategoryService {
	return &CategoryService{
		repository: repository,
	}
}

func (s *CategoryService) Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*dtos.CategoryPaginatedDto, error) {
	if categories, err := s.repository.Fetch(ctx, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var categoriesDto dtos.CategoryPaginatedDto
		for _, category := range categories.Categories {
			categoryDto := &dtos.CategoryDto{
				ID:          category.ID,
				Name:        category.Name,
				Description: category.Description,
			}

			categoriesDto.Categories = append(categoriesDto.Categories, categoryDto)
		}

		categoriesDto.TotalPage = categories.TotalPage
		categoriesDto.CurrentPage = categories.CurrentPage
		categoriesDto.NextPage = categories.NextPage
		categoriesDto.PreviousPage = categories.PreviousPage
		categoriesDto.Count = categories.Count
		categoriesDto.Size = categories.Size

		return &categoriesDto, nil
	}
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (res *dtos.CategoryDto, err error) {
	if category, err := s.repository.GetByID(ctx, id); err != nil {
		return nil, err
	} else {
		if category != nil {
			return &dtos.CategoryDto{
				ID:          category.ID,
				Name:        category.Name,
				Description: category.Description,
			}, nil
		}
		return nil, nil
	}
}

func (s *CategoryService) Update(ctx context.Context, dto *dtos.UpdateCategoryDto) error {
	return s.repository.Update(ctx, dto)
}

func (s *CategoryService) Create(ctx context.Context, dto *dtos.CreateCategoryDto) error {
	return s.repository.Create(ctx, dto)
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func (s *CategoryService) Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*dtos.CategoryPaginatedDto, error) {
	if categories, err := s.repository.Search(ctx, q, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var categoriesDto dtos.CategoryPaginatedDto
		for _, category := range categories.Categories {
			categoryDto := &dtos.CategoryDto{
				ID:          category.ID,
				Name:        category.Name,
				Description: category.Description,
			}

			categoriesDto.Categories = append(categoriesDto.Categories, categoryDto)
		}

		categoriesDto.TotalPage = categories.TotalPage
		categoriesDto.CurrentPage = categories.CurrentPage
		categoriesDto.NextPage = categories.NextPage
		categoriesDto.PreviousPage = categories.PreviousPage
		categoriesDto.Count = categories.Count
		categoriesDto.Size = categories.Size

		return &categoriesDto, nil
	}
}

func (s *CategoryService) GetProducts(ctx context.Context, id int, page int, size int, sortBy string, orderBy string) (*dtos.CategoryProductsPaginatedDto, error) {
	if products, err := s.repository.GetProducts(ctx, id, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var productsDto dtos.CategoryProductsPaginatedDto
		for _, product := range products.Products {
			productDto := &dtos.ProductDto{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				CategoryID:  product.CategoryID,
			}

			productsDto.Products = append(productsDto.Products, productDto)

			for _, image := range product.Images {
				imageDto := &dtos.ImageDto{
					ID:           image.ID,
					Name:         image.Name,
					ImageUrl:     image.ImageUrl,
					ThumbnailUrl: image.ThumbnailUrl,
				}

				productDto.Images = append(productDto.Images, imageDto)
			}
		}

		productsDto.ID = products.ID
		productsDto.Name = products.Name
		productsDto.Description = products.Description

		productsDto.TotalPage = products.TotalPage
		productsDto.CurrentPage = products.CurrentPage
		productsDto.NextPage = products.NextPage
		productsDto.PreviousPage = products.PreviousPage
		productsDto.Count = products.Count
		productsDto.Size = products.Size

		return &productsDto, nil
	}
}
