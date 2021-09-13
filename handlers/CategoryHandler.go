package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ysfada/product-management-system/domain/common"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type CategoryHandler struct {
	service interfaces.ICategoryService
}

func NewCategoryHandler(service interfaces.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

var _ interfaces.ICategoryHandler = (*CategoryHandler)(nil)

func (h *CategoryHandler) UseHandler(r fiber.Router) {
	categoriesRouter := r.Group("categories")

	categoriesRouter.Get("/", h.Fetch)
	categoriesRouter.Post("/", common.JwtMiddleware, h.Create)
	categoriesRouter.Get("/search", h.Search)
	categoriesRouter.Get("/:id/products", h.GetProducts)
	categoriesRouter.Get("/:id", h.GetByID)
	categoriesRouter.Put("/:id", common.JwtMiddleware, h.Update)
	categoriesRouter.Delete("/:id", common.JwtMiddleware, h.Delete)
}

// Category godoc
// @Summary Get categories
// @Description Get all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} dtos.CategoryPaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Router /categories [get]
func (h *CategoryHandler) Fetch(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	sortBy := strings.ToLower(c.Query("sortBy", "id"))
	if sortBy != "id" && sortBy != "name" {
		sortBy = "id"
	}
	orderBy := strings.ToUpper(c.Query("orderBy", "ASC"))
	if orderBy != "ASC" && orderBy != "DESC" {
		orderBy = "ASC"
	}

	if categories, err := h.service.Fetch(c.Context(), page, size, sortBy, orderBy); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(categories)
	}
}

// Category godoc
// @Summary Get category by id
// @Description Get category by id
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} dtos.CategoryDto
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	if id, err := c.ParamsInt("id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	} else {
		if category, err := h.service.GetByID(c.Context(), id); err != nil {
			switch err {
			case common.ErrNotFound:
				return c.SendStatus(fiber.StatusNotFound)
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(err)
			}
		} else {
			return c.JSON(category)
		}
	}
}

// Category godoc
// @Summary Update category
// @Description Update category by id
// @Tags categories
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param dto body dtos.UpdateCategoryDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	var body dtos.UpdateCategoryDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	if c.Params("id") != fmt.Sprint(body.ID) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.Update(c.Context(), &body); err != nil {
		switch err {
		case common.ErrBadParamInput:
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Category godoc
// @Summary Create category
// @Description Create new category
// @Tags categories
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.CreateCategoryDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /categories [post]
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var body dtos.CreateCategoryDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.service.Create(c.Context(), &body); err != nil {
		switch err {
		case common.ErrBadParamInput:
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Category godoc
// @Summary Delete category
// @Description Delete category by id
// @Tags categories
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param Authorization header string true "Bearer"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	if id, err := c.ParamsInt("id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	} else {
		if err := h.service.Delete(c.Context(), id); err != nil {
			switch err {
			case common.ErrNotFound:
				return c.SendStatus(fiber.StatusNotFound)
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(err)
			}
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Category godoc
// @Summary Search category
// @Description Search categories by category name
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} dtos.CategoryPaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param q query string true "query string to search in name"
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Router /categories/search [get]
func (h *CategoryHandler) Search(c *fiber.Ctx) error {
	q := c.Query("q")
	if len(q) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	sortBy := strings.ToLower(c.Query("sortBy", "id"))
	if sortBy != "id" && sortBy != "name" {
		sortBy = "id"
	}
	orderBy := strings.ToUpper(c.Query("orderBy", "ASC"))
	if orderBy != "ASC" && orderBy != "DESC" {
		orderBy = "ASC"
	}

	if categories, err := h.service.Search(c.Context(), q, page, size, sortBy, orderBy); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(categories)
	}
}

// Category godoc
// @Summary Get products
// @Description Get all products belongs to category
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ProductDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Param id path int true "id"
// @Router /categories/{id}/products [get]
func (h *CategoryHandler) GetProducts(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	sortBy := strings.ToLower(c.Query("sortBy", "id"))
	if sortBy != "id" && sortBy != "name" {
		sortBy = "id"
	}
	orderBy := strings.ToUpper(c.Query("orderBy", "ASC"))
	if orderBy != "ASC" && orderBy != "DESC" {
		orderBy = "ASC"
	}

	if products, err := h.service.GetProducts(c.Context(), id, page, size, sortBy, orderBy); err != nil {
		switch err {
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(products)
	}
}
