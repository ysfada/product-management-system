package interfaces

import (
	"github.com/gofiber/fiber/v2"
)

type IUserHandler interface {
	Signup(cc *fiber.Ctx) error
	Signin(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	ChangeUsername(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
}
