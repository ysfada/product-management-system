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

type ProductRepository struct {
	dbConn *pgxpool.Pool
}

var _ interfaces.IProductRepository = (*ProductRepository)(nil)

func NewProductRepository(dbConn *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		dbConn: dbConn,
	}
}

func (r *ProductRepository) Fetch(ctx context.Context, page int, size int, sortBy string, orderBy string) (*entities.ProductPaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
	(SELECT COUNT(*)
		FROM "public"."product" "p") "count",

	(SELECT JSONB_AGG("result".*)
		FROM
			(SELECT "p"."id",
					"p"."name",
					COALESCE("p"."description", '') "description",
					"p"."category_id",
					JSONB_BUILD_OBJECT(
                        'id', "c"."id",
                        'name', "c"."name",
                        'description', COALESCE("c"."description", ''),
                        'created_at', "c"."created_at",
                        'updated_at', "c"."updated_at",
                        'deleted_at', "c"."deleted_at"
                    ) "category",
					"p"."created_at",
					"p"."updated_at",
					"p"."deleted_at",
					"product_images"."images"
				FROM "public"."product" "p"
				JOIN "public"."category" "c" ON "c"."id" = "p"."category_id"
				CROSS JOIN LATERAL
					(SELECT JSONB_AGG("images") "images"
						FROM
							(SELECT "i"."id",
									"i"."name",
									"i"."image_url",
									"i"."thumbnail_url",
									"i"."created_at",
									"i"."updated_at",
									"i"."deleted_at"
								FROM "public"."product_images" "pi"
								JOIN "public"."image" "i" ON "i"."id" = "pi"."image_id"
								WHERE "pi"."product_id" = "p"."id"
								GROUP BY "pi"."id",
									"c"."name",
									"i"."id"
								ORDER BY "i"."id" ASC) "images") "product_images"
				ORDER BY "p"."%s" %s
				OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY) "result") "products"
        `, sortBy, orderBy)

	var products entities.ProductPaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, (page-1)*size, size).Scan(
		&products.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &products.Products); err != nil {
		return nil, err
	}

	products.Size = size
	products.TotalPage = int(math.Ceil(float64(products.Count) / float64(size)))
	products.CurrentPage = page
	if products.CurrentPage <= products.TotalPage && products.CurrentPage > 1 {
		products.PreviousPage = products.CurrentPage - 1
	} else {
		products.PreviousPage = -1
	}
	if products.CurrentPage < products.TotalPage {
		products.NextPage = products.CurrentPage + 1
	} else {
		products.NextPage = -1
	}

	return &products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*entities.Product, error) {
	sql := `
    SELECT  "p"."id",
            "p"."name",
            COALESCE("p"."description", '') "description",
            "p"."category_id",
            "c"."id",
            "c"."name",
            COALESCE("c"."description", '') "description",
            "c"."created_at",
            "c"."updated_at",
            "c"."deleted_at",
            "p"."created_at",
            "p"."updated_at",
            "p"."deleted_at",
            "product_images"."images"
    FROM "public"."product" "p"
    JOIN "public"."category" "c" ON "c"."id" = "p"."category_id"
    CROSS JOIN LATERAL
        (SELECT JSONB_AGG("images") "images"
            FROM
                (SELECT "i"."id",
                        "i"."name",
                        "i"."image_url",
                        "i"."thumbnail_url",
                        "i"."created_at",
                        "i"."updated_at",
                        "i"."deleted_at"
                    FROM "public"."product_images" "pi"
                    JOIN "public"."image" "i" ON "i"."id" = "pi"."image_id"
                    WHERE "pi"."product_id" = "p"."id"
                    GROUP BY "pi"."id", "c"."name", "i"."id"
                    ORDER BY "i"."id" ASC) "images") "product_images"
    WHERE "p"."id" = $1
    LIMIT 1
    `
	var product entities.Product
	var images json.RawMessage
	var category entities.Category
	if err := r.dbConn.QueryRow(ctx, sql, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.DeletedAt,
		&images,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	product.Category = &category

	if err := json.Unmarshal([]byte(images), &product.Images); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Search(ctx context.Context, q string, page int, size int, sortBy string, orderBy string) (*entities.ProductPaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
	(SELECT COUNT(*)
		FROM "public"."product" "p"
		WHERE "p"."name" LIKE '%%' || $1 || '%%') "count",

	(SELECT JSONB_AGG(T.*)
		FROM
			(SELECT "p"."id",
					"p"."name",
					COALESCE("p"."description", '') "description",
					"p"."category_id",
					JSONB_BUILD_OBJECT(
                        'id', "c"."id",
						'name', "c"."name",
						'description', COALESCE("c"."description", ''),
						'created_at', "c"."created_at",
						'updated_at', "c"."updated_at",
						'deleted_at', "c"."deleted_at"
                    ) "category",
					"p"."created_at",
					"p"."updated_at",
					"p"."deleted_at",
					"product_images"."images"
				FROM "public"."product" "p"
				JOIN "public"."category" "c" ON "c"."id" = "p"."category_id"
				CROSS JOIN LATERAL
					(SELECT JSONB_AGG("images") "images"
						FROM
							(SELECT "i"."id",
									"i"."name",
									"i"."image_url",
									"i"."thumbnail_url",
									"i"."created_at",
									"i"."updated_at",
									"i"."deleted_at"
								FROM "public"."product_images" "pi"
								JOIN "public"."image" "i" ON "i"."id" = "pi"."image_id"
								WHERE "pi"."product_id" = "p"."id"
								GROUP BY "pi"."id", "c"."name", "i"."id"
								ORDER BY "i"."id" ASC) "images") "product_images"
				WHERE "p"."name" LIKE '%%' || $1 || '%%'
				ORDER BY "p"."%s" %s
				OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY) AS T) AS ROWS
        `, sortBy, orderBy)

	var products entities.ProductPaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, q, (page-1)*size, size).Scan(
		&products.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &products.Products); err != nil {
		return nil, err
	}

	products.Size = size
	products.TotalPage = int(math.Ceil(float64(products.Count) / float64(size)))
	products.CurrentPage = page
	if products.CurrentPage <= products.TotalPage && products.CurrentPage > 1 {
		products.PreviousPage = products.CurrentPage - 1
	} else {
		products.PreviousPage = -1
	}
	if products.CurrentPage < products.TotalPage {
		products.NextPage = products.CurrentPage + 1
	} else {
		products.NextPage = -1
	}

	return &products, nil
}

func (r *ProductRepository) Update(ctx context.Context, dto *dtos.UpdateProductDto) error {
	sql := `
    UPDATE "public"."product"
    SET "name" = $1,
        "description" = $2,
        "category_id" = $3
    WHERE "id" = $4
    `

	_, err := r.dbConn.Exec(ctx, sql, dto.Name, dto.Description, dto.CategoryID, dto.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
	}
	return err
}

func (r *ProductRepository) Create(ctx context.Context, dto *dtos.CreateProductDto) error {
	sql := `
    INSERT INTO "public"."product" ("name", description", "category_id")
    VALUES ($1, $2, $3)
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.Name, dto.Description, dto.CategoryID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
	}
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	sql := `
    DELETE
    FROM "public"."product"
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

func (r *ProductRepository) GetImages(ctx context.Context, id int) ([]*entities.Image, error) {
	sql := `
    SELECT  "i"."id",
            COALESCE("i"."name", '') "name",
            "i"."image_url",
            "i"."thumbnail_url",
            "i"."created_at",
            "i"."updated_at",
            "i"."deleted_at"
    FROM "public"."product_images" "pi"
    JOIN "public"."image" "i" ON "i"."id" = "pi"."image_id"
    WHERE "pi"."product_id" = $1
    ORDER BY "i"."id"
    `
	var images []*entities.Image
	if rows, err := r.dbConn.Query(ctx, sql, id); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var image entities.Image
			if err := rows.Scan(
				&image.ID,
				&image.Name,
				&image.ImageUrl,
				&image.ThumbnailUrl,
				&image.CreatedAt,
				&image.UpdatedAt,
				&image.DeletedAt,
			); err != nil {
				return nil, err
			}
			images = append(images, &image)
		}
	}

	return images, nil
}

func (r *ProductRepository) AddImage(ctx context.Context, id int, imageID int) error {
	sql := `
    INSERT INTO "public"."product_images" ("product_id", "image_id")
    VALUES ($1, $2)
    `
	_, err := r.dbConn.Exec(ctx, sql, id, imageID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
	}
	return err
}

func (r *ProductRepository) RemoveImage(ctx context.Context, id int, imageID int) error {
	sql := `
    DELETE
    FROM "public"."product_images"
    WHERE "product_id" = $1
        AND "image_id" = $2
    `
	if cmd, err := r.dbConn.Exec(ctx, sql, id, imageID); err != nil {
		return err
	} else {
		if cmd.RowsAffected() > 0 {
			return nil
		} else {
			return common.ErrNotFound
		}
	}
}

func (r *ProductRepository) FetchVariants(ctx context.Context, id int, page int, size int, sortBy string, orderBy string) (*entities.ProductVariantPaginated, error) {
	sql := fmt.Sprintf(`
    SELECT
        (SELECT COUNT(*)
            FROM "public"."product_variant" "pv"
            WHERE "pv"."product_id" = $1) "count",
        JSONB_BUILD_OBJECT(
            'id', "p"."id",
            'name', "p"."name",
            'description', COALESCE("p"."description", ''),
            'category_id', "p"."category_id",
            'category', JSONB_BUILD_OBJECT(
                'id', "c"."id",
                'name', "c"."name",
                'description', COALESCE("c"."description", ''),
                'created_at', "c"."created_at",
                'updated_at', "c"."updated_at",
                'deleted_at', "c"."deleted_at"
            ),
            'created_at', "p"."created_at",
            'updated_at', "p"."updated_at",
            'deleted_at', "p"."deleted_at",
            'product_variants', "product_variants"."variants"
        ) "result"
    FROM "product" "p"
    JOIN "public"."category" "c" ON "c"."id" = "p"."category_id"
    CROSS JOIN LATERAL
        (SELECT JSONB_AGG("variants") "variants"
            FROM
                (SELECT "pv"."id",
                        "pv"."name",
                        "pv"."product_id",
                        "pv"."price",
                        "pv"."stock",
                        "pv"."created_at",
                        "pv"."updated_at",
                        "pv"."deleted_at",
                        JSONB_AGG(JSONB_BUILD_OBJECT(
                            'id', "attr"."id",
                            'name', "attr"."name",
                            'type',"attr"."type"
                        )
                ) "attributes"
                    FROM "public"."product_variant" "pv"
                    JOIN "public"."product_attributes" "pa" ON "pa"."product_variant_id" = "pv"."id"
                    JOIN "public"."attribute" "attr" ON "attr"."id" = "pa"."attribute_id"
                    WHERE "product_id" = "p"."id"
                    GROUP BY "pv"."id", "c"."name"
                    ORDER BY "pv"."%s" %s
                    OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY) "variants") "product_variants"
    WHERE "p"."id" = $1
    LIMIT 1
    `, sortBy, orderBy)

	var product_variants entities.ProductVariantPaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, id, (page-1)*size, size).Scan(
		&product_variants.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &product_variants); err != nil {
		return nil, err
	}

	product_variants.Size = size
	product_variants.TotalPage = int(math.Ceil(float64(product_variants.Count) / float64(size)))
	product_variants.CurrentPage = page
	if product_variants.CurrentPage <= product_variants.TotalPage && product_variants.CurrentPage > 1 {
		product_variants.PreviousPage = product_variants.CurrentPage - 1
	} else {
		product_variants.PreviousPage = -1
	}
	if product_variants.CurrentPage < product_variants.TotalPage {
		product_variants.NextPage = product_variants.CurrentPage + 1
	} else {
		product_variants.NextPage = -1
	}

	return &product_variants, nil
}

func (r *ProductRepository) GetVariantByID(ctx context.Context, id int, variantID int) (*entities.ProductVariant, error) {
	sql := `
    SELECT
        JSONB_BUILD_OBJECT(
            'id', "pv"."id",
            'name', "pv"."name",
            'product_id', "pv"."product_id",
            'product',
            JSONB_BUILD_OBJECT(
                'id', "p"."id",
                'name', "p"."name",
                'description', COALESCE("p"."description", ''),
                'category_id', "p"."category_id",
                'category', JSONB_BUILD_OBJECT(
                    'id', "c"."id",
                    'name', "c"."name",
                    'description', COALESCE("c"."description", ''),
                    'created_at', "c"."created_at",
                    'updated_at', "c"."updated_at",
                    'deleted_at', "c"."deleted_at"
                ),
                'created_at', "p"."created_at",
                'updated_at', "p"."updated_at",
                'deleted_at', "p"."deleted_at"
            ),
            'price', "pv"."price",
            'stock', "pv"."stock",
            'created_at', "pv"."created_at",
            'updated_at', "pv"."updated_at",
            'deleted_at', "pv"."deleted_at",
            'images', JSONB_AGG(
                JSONB_BUILD_OBJECT(
                    'id', "i"."id",
                    'name', "i"."name",
                    'image_url', "i"."image_url",
                    'thumbnail_url', "i"."thumbnail_url",
                    'created_at', "i"."created_at",
                    'updated_at', "i"."updated_at",
                    'deleted_at', "i"."deleted_at"
                )
            ), 'attributes',
            JSONB_AGG(
                JSONB_BUILD_OBJECT(
                    'id', "attr"."id",
                    'name', "attr"."name",
                    'type', "attr"."type"
                )
            )
        )
    FROM "public"."product_variant" "pv"
    JOIN "public"."product" "p" ON "p"."id" = "pv"."product_id"
    JOIN "public"."category" "c" ON "c"."id" = "p"."category_id"
    JOIN "public"."product_attributes" "pa" ON "pa"."product_variant_id" = "pv"."id"
    JOIN "public"."attribute" "attr" ON "attr"."id" = "pa"."attribute_id"
    JOIN "public"."product_images" "pi" ON "pi"."product_id" = "p"."id"
    JOIN "public"."image" "i" ON "i"."id" = "pi"."image_id"
    WHERE "pv"."id" = $2
    AND "pv"."product_id" = $1
    GROUP BY "pv"."id", "p"."id", "c"."id"
    LIMIT 1
    `

	var productVariant entities.ProductVariant
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(ctx, sql, id, variantID).Scan(&rows); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &productVariant); err != nil {
		return nil, err
	}

	return &productVariant, nil
}

func (r *ProductRepository) CreateVariant(ctx context.Context, dto *dtos.CreateProductVariantDto) error {
	sql := `
    INSERT INTO "public"."product_variant" ("product_id", "name", "price", "stock")
    VALUES ($1, $2, $3, $4)
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.ProductId, dto.Name, dto.Price, dto.Stock)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
		switch pgErr.Code {
		case pgerrcode.CheckViolation:
			return common.ErrBadParamInput
		default:
			return err
		}
	}
	return err
}

func (r *ProductRepository) UpdateVariant(ctx context.Context, dto *dtos.UpdateProductVariantDto) error {
	sql := `
    UPDATE "public"."product_variant"
    SET "product_id" = $1,
        "name" = $2,
        "price" = $3,
        "stock" = $4
    WHERE "id" = $5
    `

	_, err := r.dbConn.Exec(ctx, sql, dto.ProductId, dto.Name, dto.Price, dto.Stock, dto.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.CheckViolation:
			return common.ErrBadParamInput
		default:
			return err
		}
	}
	return err
}

func (r *ProductRepository) DeleteVariant(ctx context.Context, id int, variantID int) error {
	sql := `
    DELETE
    FROM "public"."product_variant"
    WHERE "id" = $2 AND "product_id" = $1
    `
	if cmd, err := r.dbConn.Exec(ctx, sql, id, variantID); err != nil {
		return err
	} else {
		if cmd.RowsAffected() > 0 {
			return nil
		} else {
			return common.ErrNotFound
		}
	}
}

func (r *ProductRepository) GetAttributes(ctx context.Context, id int, variantID int) ([]*entities.Attribute, error) {
	sql := `
    SELECT  "a"."id",
            "a"."name",
            "a"."type",
            "a"."created_at",
            "a"."updated_at",
            "a"."deleted_at"
    FROM "public"."product_attributes" "pa"
    JOIN "public"."attribute" "a" ON "a"."id" = "pa"."attribute_id"
    WHERE "pa"."product_variant_id" = $1
    ORDER BY "a"."id"
    `
	var attributes []*entities.Attribute
	if rows, err := r.dbConn.Query(ctx, sql, variantID); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var attribute entities.Attribute
			if err := rows.Scan(
				&attribute.ID,
				&attribute.Name,
				&attribute.Type,
				&attribute.CreatedAt,
				&attribute.UpdatedAt,
				&attribute.DeletedAt,
			); err != nil {
				return nil, err
			}
			attributes = append(attributes, &attribute)
		}
	}

	return attributes, nil
}

func (r *ProductRepository) AddAttribute(ctx context.Context, dto *dtos.CreateProductVariantAttributeDto) error {
	sql := `
    INSERT INTO "public"."product_attributes" ("product_variant_id", "attribute_id")
    VALUES ($1, $2)
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.ProductVariantID, dto.AttributeID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
		switch pgErr.Code {
		case pgerrcode.CheckViolation:
			return common.ErrBadParamInput
		default:
			return err
		}
	}
	return err
}

func (r *ProductRepository) RemoveAttribute(ctx context.Context, id int, variantID int, attributeID int) error {
	sql := `
    DELETE
    FROM "public"."product_attributes"
    WHERE "product_variant_id" = $1 AND "attribute_id" = $2
    `
	if cmd, err := r.dbConn.Exec(ctx, sql, variantID, attributeID); err != nil {
		return err
	} else {
		if cmd.RowsAffected() > 0 {
			return nil
		} else {
			return common.ErrNotFound
		}
	}
}

func (r *ProductRepository) SearchVariants(ctx context.Context, q string, id int, page int, size int, sortBy string, orderBy string, attrs []*dtos.AttributeSearchQueryDto) (*entities.ProductVariantPaginated, error) {
	base := `
	AND (EXISTS
		(SELECT 1
			FROM "public"."product_attributes" "pa"
			JOIN "public"."attribute" "a" ON "a"."id" = "pa"."attribute_id"
			WHERE "pa"."product_variant_id" = "pv"."id"
			GROUP BY "a"."id"
			HAVING "a"."type" = %s AND "a"."name" IN (%s)
		)
	)
	`
	baseFilled := ""

	j := 5
	for i := 0; i < len(attrs); i++ {
		n := ""
		for ii := 0; ii < len(attrs[i].Names); ii++ {
			n += fmt.Sprintf("$%d,", j)
			j += 1
		}
		n = n[:len(n)-1]
		baseFilled += fmt.Sprintf(base, fmt.Sprintf("$%d", j), n)
		j += 1
	}

	sql := fmt.Sprintf(`
    SELECT
        (SELECT COUNT(*)
            FROM "public"."product_variant" "pv"
            WHERE "pv"."product_id" = $1
                    AND "pv"."name" LIKE '%%' || $4 || '%%'
                    %s
        ) "count",
        JSONB_BUILD_OBJECT(
            'id', "p"."id",
            'name', "p"."name",
            'description', COALESCE("p"."description", ''),
            'category_id', "p"."category_id",
            'category', JSONB_BUILD_OBJECT(
                'id', "c"."id",
                'name', "c"."name",
                'description', COALESCE("c"."description", ''),
                'created_at', "c"."created_at",
                'updated_at', "c"."updated_at",
                'deleted_at',"c"."deleted_at"
            ),
            'created_at', "p"."created_at",
            'updated_at', "p"."updated_at",
            'deleted_at', "p"."deleted_at",
            'product_variants', "product_variants"."variants"
        ) "result"
    FROM "product" "p"
    JOIN "public"."category" "c" ON "c"."id" = "p"."category_id"
    CROSS JOIN LATERAL
        (SELECT JSONB_AGG("variants") "variants"
            FROM
                (SELECT "pv"."id",
                        "pv"."name",
                        "pv"."product_id",
                        "pv"."price",
                        "pv"."stock",
                        "pv"."created_at",
                        "pv"."updated_at",
                        "pv"."deleted_at",
                        JSONB_AGG(
                            JSONB_BUILD_OBJECT(
                                'id', "attr"."id",
                                'name', "attr"."name",
                                'type', "attr"."type"
                            )
                        ) "attributes"
                FROM "public"."product_variant" "pv"
                JOIN "public"."product_attributes" "pa" ON "pa"."product_variant_id" = "pv"."id"
                JOIN "public"."attribute" "attr" ON "attr"."id" = "pa"."attribute_id"
                WHERE "pv"."product_id" = "p"."id"
                    AND "pv"."name" LIKE '%%' || $4 || '%%'
                    %s
                GROUP BY "pv"."id"
                ORDER BY "pv"."%s" %s
                OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY) "variants") "product_variants"
    WHERE "p"."id" = $1
    LIMIT 1
    `, baseFilled, baseFilled, sortBy, orderBy)

	args := []interface{}{
		id,
		(page - 1) * size,
		size,
		q,
	}
	for _, attr := range attrs {
		for _, name := range attr.Names {
			args = append(args, name)
		}
		args = append(args, attr.Type)
	}

	var product_variants entities.ProductVariantPaginated
	var rows json.RawMessage
	if err := r.dbConn.QueryRow(
		ctx,
		sql,
		args...,
	).Scan(
		&product_variants.Count,
		&rows,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(rows), &product_variants); err != nil {
		return nil, err
	}

	product_variants.Size = size
	product_variants.TotalPage = int(math.Ceil(float64(product_variants.Count) / float64(size)))
	product_variants.CurrentPage = page
	if product_variants.CurrentPage <= product_variants.TotalPage && product_variants.CurrentPage > 1 {
		product_variants.PreviousPage = product_variants.CurrentPage - 1
	} else {
		product_variants.PreviousPage = -1
	}
	if product_variants.CurrentPage < product_variants.TotalPage {
		product_variants.NextPage = product_variants.CurrentPage + 1
	} else {
		product_variants.NextPage = -1
	}

	return &product_variants, nil
}
