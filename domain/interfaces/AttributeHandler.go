package interfaces

import "github.com/gofiber/fiber/v2"

type IAttributeHandler interface {
	Fetch(cc *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
}
