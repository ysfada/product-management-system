package interfaces

import "github.com/gofiber/fiber/v2"

type IImageHandler interface {
	Add(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}
