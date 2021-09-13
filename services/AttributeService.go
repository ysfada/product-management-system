package services

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type AttributeService struct {
	repository interfaces.IAttributeRepository
}

var _ interfaces.IAttributeService = (*AttributeService)(nil)

func NewAttributeService(repository interfaces.IAttributeRepository) *AttributeService {
	return &AttributeService{
		repository: repository,
	}
}

func (s *AttributeService) Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*dtos.AttributePaginatedDto, error) {
	if attributes, err := s.repository.Fetch(ctx, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var attributesDto dtos.AttributePaginatedDto
		for _, attribute := range attributes.Attributes {
			attributeDto := &dtos.AttributeDto{
				ID:   attribute.ID,
				Name: attribute.Name,
				Type: attribute.Type,
			}

			attributesDto.Attributes = append(attributesDto.Attributes, attributeDto)
		}

		attributesDto.TotalPage = attributes.TotalPage
		attributesDto.CurrentPage = attributes.CurrentPage
		attributesDto.NextPage = attributes.NextPage
		attributesDto.PreviousPage = attributes.PreviousPage
		attributesDto.Count = attributes.Count
		attributesDto.Size = attributes.Size

		return &attributesDto, nil
	}
}

func (s *AttributeService) GetByID(ctx context.Context, id int) (res *dtos.AttributeDto, err error) {
	if attribute, err := s.repository.GetByID(ctx, id); err != nil {
		return nil, err
	} else {
		if attribute != nil {
			return &dtos.AttributeDto{
				ID:   attribute.ID,
				Name: attribute.Name,
				Type: attribute.Type,
			}, nil
		}
		return nil, nil
	}
}

func (s *AttributeService) Update(ctx context.Context, dto *dtos.UpdateAttributeDto) error {
	return s.repository.Update(ctx, dto)
}

func (s *AttributeService) Create(ctx context.Context, dto *dtos.CreateAttributeDto) error {
	return s.repository.Create(ctx, dto)
}

func (s *AttributeService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func (s *AttributeService) Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*dtos.AttributePaginatedDto, error) {
	if attributes, err := s.repository.Search(ctx, q, page, size, sortBy, orderBy); err != nil {
		return nil, err
	} else {
		var attributesDto dtos.AttributePaginatedDto
		for _, category := range attributes.Attributes {
			categoryDto := &dtos.AttributeDto{
				ID:   category.ID,
				Name: category.Name,
				Type: category.Type,
			}

			attributesDto.Attributes = append(attributesDto.Attributes, categoryDto)
		}

		attributesDto.TotalPage = attributes.TotalPage
		attributesDto.CurrentPage = attributes.CurrentPage
		attributesDto.NextPage = attributes.NextPage
		attributesDto.PreviousPage = attributes.PreviousPage
		attributesDto.Count = attributes.Count
		attributesDto.Size = attributes.Size

		return &attributesDto, nil
	}
}
