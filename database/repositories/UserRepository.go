package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ysfada/product-management-system/domain/common"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/entities"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type UserRepository struct {
	dbConn *pgxpool.Pool
}

var _ interfaces.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(dbConn *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		dbConn: dbConn,
	}
}

func (r *UserRepository) Create(ctx context.Context, dto *dtos.SignupDto) error {
	sql := `
        INSERT INTO "public"."user" ("username", "password")
        VALUES ($1, $2)
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.Username, dto.Password)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
		if pgErr.Code == pgerrcode.UniqueViolation {
			return common.ErrConflict
		}
	}
	return err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (res *entities.User, err error) {
	sql := `
    SELECT  "u"."id",
            "u"."username",
            "u"."password",
            "u"."is_active",
            "u"."is_staff",
            "u"."is_superuser",
            "u"."created_at",
            "u"."updated_at",
            "u"."deleted_at"
    FROM "public"."user" "u"
    WHERE "u"."username" = $1
    LIMIT 1
    `
	var user entities.User
	if err := r.dbConn.QueryRow(ctx, sql, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&user.IsStaff,
		&user.IsSuperuser,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, common.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, username string) error {
	sql := `
    DELETE
    FROM "public"."user"
    WHERE "username" = $1
    `
	if cmd, err := r.dbConn.Exec(ctx, sql, username); err != nil {
		return err
	} else {
		if cmd.RowsAffected() > 0 {
			return nil
		} else {
			return common.ErrNotFound
		}
	}
}

func (r *UserRepository) ChangeUsername(ctx context.Context, username string, dto *dtos.ChangeUsernameDto) error {
	sql := `
    UPDATE "public"."user"
    SET "username" = $1
    WHERE "username" = $2
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.Username, username)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
		if pgErr.Code == pgerrcode.UniqueViolation {
			return common.ErrConflict
		}
	}
	return err
}

func (r *UserRepository) ChangePassword(ctx context.Context, username string, dto *dtos.ChangePasswordDto) error {
	sql := `
    UPDATE "public"."user"
    SET "password" = $1
    WHERE "username" = $2
    `
	_, err := r.dbConn.Exec(ctx, sql, dto.Password, username)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.CheckViolation {
			return common.ErrBadParamInput
		}
		if pgErr.Code == pgerrcode.UniqueViolation {
			return common.ErrConflict
		}
	}
	return err
}
