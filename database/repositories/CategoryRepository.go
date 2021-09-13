package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ysfada/product-management-system/domain/common"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/entities"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type CategoryRepository struct {
	dbConn *pgxpool.Pool
}

var _ interfaces.ICategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(dbConn *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		dbConn: dbConn,
	}
}

func (r *CategoryRepository) Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*entities.CategoryPaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
	(SELECT COUNT(*)
		FROM "public"."category" "c") "count",

	(SELECT JSONB_AGG("result".*)
		FROM
			(SELECT "c"."id",
					"c"."name",
					COALESCE("c"."description", '') "description",
					"c"."created_at",
					"c"."updated_at",
					"c"."deleted_at"
				FROM "public"."category" "c"
				ORDER BY "c"."%s" %s
				OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY) "result") "categories"
        `, sortBy, orderBy)

	var categories entities.CategoryPaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, (page-1)*size, size).Scan(
		&categories.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &categories.Categories); err != nil {
		return nil, err
	}

	categories.Size = size
	categories.TotalPage = int(math.Ceil(float64(categories.Count) / float64(size)))
	categories.CurrentPage = page
	if categories.CurrentPage <= categories.TotalPage && categories.CurrentPage > 1 {
		categories.PreviousPage = categories.CurrentPage - 1
	} else {
		categories.PreviousPage = -1
	}
	if categories.CurrentPage < categories.TotalPage {
		categories.NextPage = categories.CurrentPage + 1
	} else {
		categories.NextPage = -1
	}

	return &categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (res *entities.Category, err error) {
	sql := `
    SELECT  "c"."id",
	        "c"."name",
	        COALESCE("c"."description", '') AS DESCRIPTION,
	        "c"."created_at",
	        "c"."updated_at",
	        "c"."deleted_at"
    FROM "public"."category" "c"
    WHERE "id" = $1
    LIMIT 1
    `
	var category entities.Category
	if err := r.dbConn.QueryRow(ctx, sql, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (r *CategoryRepository) Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*entities.CategoryPaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
    (SELECT COUNT(*)
        FROM "public"."category" "c"
        WHERE "c"."name" LIKE '%%' || $1 || '%%') "count",

    (SELECT JSONB_AGG("result".*)
        FROM (
            SELECT "c"."id",
                "c"."name",
                COALESCE("c"."description", '') AS DESCRIPTION,
                "c"."created_at",
                "c"."updated_at",
                "c"."deleted_at"
            FROM "public"."category" "c"
            WHERE "c"."name" LIKE '%%' || $1 || '%%'
            ORDER BY "c"."%s" %s
            OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY) "result") "categories"
        `, sortBy, orderBy)

	var categories entities.CategoryPaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, q, (page-1)*size, size).Scan(
		&categories.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &categories.Categories); err != nil {
		return nil, err
	}

	categories.Size = size
	categories.TotalPage = int(math.Ceil(float64(categories.Count) / float64(size)))
	categories.CurrentPage = page
	if categories.CurrentPage <= categories.TotalPage && categories.CurrentPage > 1 {
		categories.PreviousPage = categories.CurrentPage - 1
	} else {
		categories.PreviousPage = -1
	}
	if categories.CurrentPage < categories.TotalPage {
		categories.NextPage = categories.CurrentPage + 1
	} else {
		categories.NextPage = -1
	}

	return &categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, dto *dtos.UpdateCategoryDto) error {
	sql := `
    UPDATE "public"."category"
    SET "name" = $1,
        "description" = $2
    WHERE "id" = $3
    `

	_, err := r.dbConn.Exec(ctx, sql, dto.Name, dto.Description, dto.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
	}
	return err
}

func (r *CategoryRepository) Create(ctx context.Context, dto *dtos.CreateCategoryDto) error {
	sql := `
    INSERT INTO "public"."category"("name", "description")
    VALUES ($1, $2)
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.Name, dto.Description)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
	}
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	sql := `
    DELETE
    FROM "public"."category"
    WHERE "id" = $1
    `
	if cmd, err := r.dbConn.Exec(ctx, sql, id); err != nil {
		return err
	} else {
		if cmd.RowsAffected() > 0 {
			return nil
		} else {
			return common.ErrNotFound
		}
	}
}

func (r *CategoryRepository) GetProducts(ctx context.Context, id int, page int, size int, sortBy string, orderBy string) (*entities.CategoryProductsPaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
        (SELECT COUNT(*) FROM "public"."product" "p" WHERE "p"."category_id"=$1) "count",

        jsonb_build_object(
            'id', "c"."id",
            'name', "c"."name",
            'description', COALESCE("c"."description", ''),
            'created_at', "c"."created_at",
            'updated_at', "c"."updated_at",
            'deleted_at', "c"."deleted_at",
            'products', "p"."products"
        ) "products"
        FROM "category" "c"
        CROSS JOIN LATERAL (
        SELECT jsonb_agg("products") "products"
        FROM (
            SELECT "p"."id",
                    "p"."name",
                    COALESCE("c"."description", '') "description",
                    "p"."category_id",
                    "p"."created_at",
                    "p"."updated_at",
                    "p"."deleted_at",
                    "product_images"."images"
            FROM "public"."product" "p"
            CROSS JOIN LATERAL (
                SELECT jsonb_agg("images") "images"
                FROM (
                    SELECT  "i"."id",
                            "i"."name",
                            "i"."image_url",
                            "i"."thumbnail_url",
                            "i"."created_at",
                            "i"."updated_at",
                            "i"."deleted_at"
                    FROM "public"."product_images" "pi"
                    JOIN "public"."image" "i" ON "i"."id"="pi"."image_id"
                    WHERE "pi"."product_id" = "p"."id"
                    GROUP BY "pi"."id", "c"."name", "i"."id"
                    ORDER BY "i"."id" ASC
                    ) "images"
                ) "product_images"
            WHERE "category_id" = "c"."id"
            ORDER BY "p"."%s" %s
            OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY
            ) "products"
        ) "p"
        WHERE "c"."id"=$1
        LIMIT 1
    `, sortBy, orderBy)
	var categoryProducts entities.CategoryProductsPaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, id, (page-1)*size, size).Scan(
		&categoryProducts.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &categoryProducts); err != nil {
		return nil, err
	}

	categoryProducts.Size = size
	categoryProducts.TotalPage = int(math.Ceil(float64(categoryProducts.Count) / float64(size)))
	categoryProducts.CurrentPage = page
	if categoryProducts.CurrentPage <= categoryProducts.TotalPage && categoryProducts.CurrentPage > 1 {
		categoryProducts.PreviousPage = categoryProducts.CurrentPage - 1
	} else {
		categoryProducts.PreviousPage = -1
	}
	if categoryProducts.CurrentPage < categoryProducts.TotalPage {
		categoryProducts.NextPage = categoryProducts.CurrentPage + 1
	} else {
		categoryProducts.NextPage = -1
	}

	return &categoryProducts, nil
}
