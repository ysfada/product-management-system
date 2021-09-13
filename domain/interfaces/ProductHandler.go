package interfaces

import "github.com/gofiber/fiber/v2"

type IProductHandler interface {
	Fetch(cc *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	GetImages(c *fiber.Ctx) error
	AddImage(c *fiber.Ctx) error
	RemoveImage(c *fiber.Ctx) error
	FetchVariants(c *fiber.Ctx) error
	SearchVariants(c *fiber.Ctx) error
	GetVariantByID(c *fiber.Ctx) error
	CreateVariant(c *fiber.Ctx) error
	UpdateVariant(c *fiber.Ctx) error
	DeleteVariant(c *fiber.Ctx) error
	GetAttributes(c *fiber.Ctx) error
	AddAttribute(c *fiber.Ctx) error
	RemoveAttribute(c *fiber.Ctx) error
}
