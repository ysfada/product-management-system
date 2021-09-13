package interfaces

import (
	"context"

	"github.com/ysfada/product-management-system/domain/dtos"
)

type IUserService interface {
	Signup(ctx context.Context, dto *dtos.SignupDto) error
	Signin(ctx context.Context, tdo *dtos.SigninDto) (token string, err error)
	Me(ctx context.Context, username string) (res *dtos.UserDto, err error)
	Delete(ctx context.Context, username string) error
	ChangeUsername(ctx context.Context, username string, dto *dtos.ChangeUsernameDto) error
	ChangePassword(ctx context.Context, username string, dto *dtos.ChangePasswordDto) error
}
