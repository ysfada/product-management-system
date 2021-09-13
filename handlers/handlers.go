package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ysfada/product-management-system/database"
	"github.com/ysfada/product-management-system/database/repositories"
	"github.com/ysfada/product-management-system/services"
	"github.com/ysfada/product-management-system/util/hasher"
)

func Use(r fiber.Router) {
	argon2 := hasher.NewArgon2()

	userRepository := repositories.NewUserRepository(database.DbConn)
	categoryRepository := repositories.NewCategoryRepository(database.DbConn)
	productRepository := repositories.NewProductRepository(database.DbConn)
	attributeRepository := repositories.NewAttributeRepository(database.DbConn)
	imageRepository := repositories.NewImageRepository(database.DbConn)

	userService := services.NewUserService(userRepository, argon2)
	categoryService := services.NewCategoryService(categoryRepository)
	attributeService := services.NewAttributeService(attributeRepository)
	imageService := services.NewImageService(imageRepository)
	productService := services.NewProductService(productRepository, imageService)

	NewUserHandler(userService).UseHandler(r)
	NewCategoryHandler(categoryService).UseHandler(r)
	NewProductHandler(productService).UseHandler(r)
	NewAttributeHandler(attributeService).UseHandler(r)
}
