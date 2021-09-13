package interfaces

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/entities"
)

type IUserRepository interface {
	GetByUsername(ctx context.Context, username string) (res *entities.User, err error)
	Create(ctx context.Context, dto *dtos.SignupDto) error
	Delete(ctx context.Context, username string) error
	ChangeUsername(ctx context.Context, username string, dto *dtos.ChangeUsernameDto) error
	ChangePassword(ctx context.Context, username string, dto *dtos.ChangePasswordDto) error
}
