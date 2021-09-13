package services

import (
	"context"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ysfada/product-management-system/domain/common"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
	"github.com/ysfada/product-management-system/util/hasher"
)

type UserService struct {
	repository interfaces.IUserRepository
	hasher     hasher.Hasher
	validate   *validator.Validate
}

var _ interfaces.IUserService = (*UserService)(nil)

func NewUserService(repository interfaces.IUserRepository, hasher hasher.Hasher) *UserService {
	return &UserService{
		repository: repository,
		hasher:     hasher,
		validate:   validator.New(),
	}
}

func (s *UserService) Signup(ctx context.Context, dto *dtos.SignupDto) error {
	if err := s.validate.Struct(dto); err != nil {
		return err
	}

	if hash, err := s.hasher.Hash(dto.Password); err != nil {
		return err
	} else {
		dto.Password = string(hash)
		return s.repository.Create(ctx, dto)
	}
}

func (s *UserService) Signin(ctx context.Context, dto *dtos.SigninDto) (token string, err error) {
	if err := s.validate.Struct(dto); err != nil {
		return "", err
	}

	if user, err := s.repository.GetByUsername(ctx, dto.Username); err != nil {
		return "", err
	} else {
		if user != nil {
			if err := s.hasher.Compare(user.Password, dto.Password); err != nil {
				return "", common.ErrBadParamInput
			}

			// Create token
			token := jwt.New(jwt.SigningMethodHS256)

			// Set claims
			claims := token.Claims.(jwt.MapClaims)
			claims["username"] = user.Username
			claims["isStaf"] = user.IsStaff
			claims["isSuperuser"] = user.IsSuperuser
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

			// Generate encoded token and send it as response.
			return token.SignedString([]byte(os.Getenv("SECRET")))
		}
		return "", nil
	}
}

func (s *UserService) Me(ctx context.Context, username string) (res *dtos.UserDto, err error) {
	if user, err := s.repository.GetByUsername(ctx, username); err != nil {
		return nil, err
	} else {
		if user != nil {
			return &dtos.UserDto{
				ID:          user.ID,
				Username:    user.Username,
				IsActive:    user.IsActive,
				IsStaff:     user.IsStaff,
				IsSuperuser: user.IsSuperuser,
			}, nil
		}
		return nil, nil
	}
}

func (s *UserService) Delete(ctx context.Context, username string) error {
	return s.repository.Delete(ctx, username)
}

func (s *UserService) ChangeUsername(ctx context.Context, username string, dto *dtos.ChangeUsernameDto) error {
	if err := s.validate.Struct(dto); err != nil {
		return err
	}

	return s.repository.ChangeUsername(ctx, username, dto)
}

func (s *UserService) ChangePassword(ctx context.Context, username string, dto *dtos.ChangePasswordDto) error {
	if err := s.validate.Struct(dto); err != nil {
		return err
	}

	if hash, err := s.hasher.Hash(dto.Password); err != nil {
		return err
	} else {
		dto.Password = string(hash)
		return s.repository.ChangePassword(ctx, username, dto)
	}
}
