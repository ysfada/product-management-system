package handlers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ysfada/product-management-system/domain/common"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type UserHandler struct {
	service interfaces.IUserService
}

func NewUserHandler(service interfaces.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

var _ interfaces.IUserHandler = (*UserHandler)(nil)

func (h *UserHandler) UseHandler(r fiber.Router) {
	userRouter := r.Group("users")
	userRouter.Post("/signup", h.Signup)
	userRouter.Post("/signin", h.Signin)

	userRouter.Use("/me", common.JwtMiddleware)

	userRouter.Get("/me", h.Me)
	userRouter.Delete("/me", h.Delete)
	userRouter.Post("/me/change-username", h.ChangeUsername)
	userRouter.Post("/me/change-password", h.ChangePassword)
}

// User godoc
// @Summary Signup
// @Description Create new account
// @Tags users
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.SignupDto true "dto"
// @Router /users/signup [post]
func (h *UserHandler) Signup(c *fiber.Ctx) error {
	var body dtos.SignupDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.service.Signup(c.Context(), &body); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrors.Error())
		}
		switch err {
		case common.ErrBadParamInput:
			return c.SendStatus(fiber.StatusBadRequest)
		case common.ErrConflict:
			return c.Status(fiber.StatusBadRequest).JSON("username already exists")
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// User godoc
// @Summary Signin
// @Description Signin with username and password
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.SigninDto true "dto"
// @Router /users/signin [post]
func (h *UserHandler) Signin(c *fiber.Ctx) error {
	var body dtos.SigninDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if jwt, err := h.service.Signin(c.Context(), &body); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrors.Error())
		}
		switch err {
		case common.ErrBadParamInput:
			return c.SendStatus(fiber.StatusBadRequest)
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.Status(fiber.StatusOK).JSON(jwt)
	}
}

// User godoc
// @Summary Get current users details
// @Description Get current users details
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dtos.UserDto
// @Failure 400 {object} string
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param Authorization header string true "Bearer"
// @Router /users/me [get]
func (h *UserHandler) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	username, ok := claims["username"].(string)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if user, err := h.service.Me(c.Context(), username); err != nil {
		switch err {
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(user)
	}
}

// User godoc
// @Summary Delete current user
// @Description Delete current user
// @Tags users
// @Accept json
// @Produce json
// @Success 204
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param Authorization header string true "Bearer"
// @Router /users/me [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	username, ok := claims["username"].(string)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := h.service.Delete(c.Context(), username); err != nil {
		switch err {
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// User godoc
// @Summary Change username
// @Description Change username
// @Tags users
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.ChangeUsernameDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /users/me/change-username [post]
func (h *UserHandler) ChangeUsername(c *fiber.Ctx) error {
	var body dtos.ChangeUsernameDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	username, ok := claims["username"].(string)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := h.service.ChangeUsername(c.Context(), username, &body); err != nil {
		switch err {
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// User godoc
// @Summary Change password
// @Description Change password
// @Tags users
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.ChangePasswordDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /users/me/change-password [post]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	var body dtos.ChangePasswordDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}
	username, ok := claims["username"].(string)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := h.service.ChangePassword(c.Context(), username, &body); err != nil {
		switch err {
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}
