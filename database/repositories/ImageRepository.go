package repositories

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ysfada/product-management-system/domain/common"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type ImageRepository struct {
	dbConn *pgxpool.Pool
}

var _ interfaces.IImageRepository = (*ImageRepository)(nil)

func NewImageRepository(dbConn *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{
		dbConn: dbConn,
	}
}

func (r *ImageRepository) Update(ctx context.Context, dto *dtos.UpdateImageDto) error {
	return nil
}

func (r *ImageRepository) Create(ctx context.Context, dto *dtos.CreateImageDto) (int, error) {
	sql := `
    INSERT INTO "public"."image" ("name", "image_url", "thumbnail_url")
    VALUES ($1, $2, $3) RETURNING "id"
    `
	var id int
	if err := r.dbConn.QueryRow(ctx, sql, dto.Name, dto.ImageUrl, dto.ThumbnailUrl).Scan(&id); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return -1, common.ErrNotFound
		default:
			return -1, err
		}
	}

	return id, nil
}

func (r *ImageRepository) Delete(ctx context.Context, id int) (string, error) {
	sql := `
    DELETE
    FROM "public"."image"
    WHERE "id" = $1 RETURNING "image_url"
    `
	var imageURL string
	if err := r.dbConn.QueryRow(ctx, sql, id).Scan(&imageURL); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return "", common.ErrNotFound
		default:
			return "", err
		}
	}

	return imageURL, nil
}
