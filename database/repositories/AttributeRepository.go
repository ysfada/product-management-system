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

type AttributeRepository struct {
	dbConn *pgxpool.Pool
}

var _ interfaces.IAttributeRepository = (*AttributeRepository)(nil)

func NewAttributeRepository(dbConn *pgxpool.Pool) *AttributeRepository {
	return &AttributeRepository{
		dbConn: dbConn,
	}
}

func (r *AttributeRepository) Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*entities.AttributePaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
	(SELECT COUNT(*)
		FROM "public"."attribute" "a") "count",

	(SELECT JSONB_AGG("result".*)
		FROM
			(SELECT "a"."id",
					"a"."name",
					"a"."type",
					"a"."created_at",
					"a"."updated_at",
					"a"."deleted_at"
				FROM "public"."attribute" "a"
				ORDER BY "a"."%s" %s
				OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY) "result") "attributes"
        `, sortBy, orderBy)

	var attributes entities.AttributePaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, (page-1)*size, size).Scan(
		&attributes.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &attributes.Attributes); err != nil {
		return nil, err
	}

	attributes.Size = size
	attributes.TotalPage = int(math.Ceil(float64(attributes.Count) / float64(size)))
	attributes.CurrentPage = page
	if attributes.CurrentPage <= attributes.TotalPage && attributes.CurrentPage > 1 {
		attributes.PreviousPage = attributes.CurrentPage - 1
	} else {
		attributes.PreviousPage = -1
	}
	if attributes.CurrentPage < attributes.TotalPage {
		attributes.NextPage = attributes.CurrentPage + 1
	} else {
		attributes.NextPage = -1
	}

	return &attributes, nil
}

func (r *AttributeRepository) GetByID(ctx context.Context, id int) (res *entities.Attribute, err error) {
	sql := `
    SELECT "a"."id",
        "a"."name",
        "a"."type",
        "a"."created_at",
        "a"."updated_at",
        "a"."deleted_at"
    FROM "public"."attribute" "a"
    WHERE "id" = $1
    LIMIT 1
    `
	var attribute entities.Attribute
	if err := r.dbConn.QueryRow(ctx, sql, id).Scan(
		&attribute.ID,
		&attribute.Name,
		&attribute.Type,
		&attribute.CreatedAt,
		&attribute.UpdatedAt,
		&attribute.DeletedAt,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	return &attribute, nil
}

func (r *AttributeRepository) Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*entities.AttributePaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
	(SELECT COUNT(*)
		FROM "public"."attribute" "a"
		WHERE "a"."name" LIKE '%%' || $1 || '%%') "count",

	(SELECT JSONB_AGG(result.*)
		FROM
			(SELECT "a"."id",
					"a"."name",
					"a"."type",
					"a"."created_at",
					"a"."updated_at",
					"a"."deleted_at"
				FROM "public"."attribute" "a"
				WHERE "a"."name" LIKE '%%' || $1 || '%%'
				ORDER BY "a"."%s" %s
				OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY) "result") "attributes"
        `, sortBy, orderBy)

	var attributes entities.AttributePaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, q, (page-1)*size, size).Scan(
		&attributes.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &attributes.Attributes); err != nil {
		return nil, err
	}

	attributes.Size = size
	attributes.TotalPage = int(math.Ceil(float64(attributes.Count) / float64(size)))
	attributes.CurrentPage = page
	if attributes.CurrentPage <= attributes.TotalPage && attributes.CurrentPage > 1 {
		attributes.PreviousPage = attributes.CurrentPage - 1
	} else {
		attributes.PreviousPage = -1
	}
	if attributes.CurrentPage < attributes.TotalPage {
		attributes.NextPage = attributes.CurrentPage + 1
	} else {
		attributes.NextPage = -1
	}

	return &attributes, nil
}

func (r *AttributeRepository) Update(ctx context.Context, dto *dtos.UpdateAttributeDto) error {
	sql := `
    UPDATE "public"."attribute"
    SET "name" = $1,
        "type" = $2
    WHERE "id" = $3
    `

	_, err := r.dbConn.Exec(ctx, sql, dto.Name, dto.Type, dto.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
	}
	return err
}

func (r *AttributeRepository) Create(ctx context.Context, dto *dtos.CreateAttributeDto) error {
	sql := `
    INSERT INTO "public"."attribute"("name","type")
    VALUES ($1, $2)
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.Name, dto.Type)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
	}
	return err
}

func (r *AttributeRepository) Delete(ctx context.Context, id int) error {
	sql := `
    DELETE
    FROM "public"."attribute"
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
