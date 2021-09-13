package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type ImageService struct {
	repository interfaces.IImageRepository
}

var _ interfaces.IImageService = (*ImageService)(nil)

func NewImageService(repository interfaces.IImageRepository) *ImageService {
	return &ImageService{
		repository: repository,
	}
}

func (s *ImageService) Save(ctx context.Context, id int, fileheader *multipart.FileHeader) (int, error) {
	// generate new uuid for image name
	uniqueId := uuid.New()

	// remove "- from imageName"
	filename := strings.Replace(uniqueId.String(), "-", "", -1)

	// extract image extension from original file filename
	fileExt := strings.Split(fileheader.Filename, ".")[1]

	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	// save image to ./public/images dir
	if err := fasthttp.SaveMultipartFile(fileheader, fmt.Sprintf("./public/images/%s", image)); err != nil {
		return -1, err
	}

	// generate image url to serve to client
	imageUrl := fmt.Sprintf("/public/images/%s", image)
	thumbnailUrl := fmt.Sprintf("/public/images/%s", image)

	return s.repository.Create(ctx, &dtos.CreateImageDto{
		Name:         fileheader.Filename,
		ImageUrl:     imageUrl,
		ThumbnailUrl: thumbnailUrl,
	})
}

func (s *ImageService) Remove(ctx context.Context, imageID int) error {
	if imageURL, err := s.repository.Delete(ctx, imageID); err != nil {
		return err
	} else {
		return os.Remove(fmt.Sprintf(".%s", imageURL))
	}
}
