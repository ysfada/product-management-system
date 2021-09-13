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

type AttributeHandler struct {
	service interfaces.IAttributeService
}

func NewAttributeHandler(service interfaces.IAttributeService) *AttributeHandler {
	return &AttributeHandler{
		service: service,
	}
}

var _ interfaces.IAttributeHandler = (*AttributeHandler)(nil)

func (h *AttributeHandler) UseHandler(r fiber.Router) {
	attributesRouter := r.Group("attributes")

	attributesRouter.Get("/", h.Fetch)
	attributesRouter.Post("/", common.JwtMiddleware, h.Create)
	attributesRouter.Get("/search", h.Search)
	attributesRouter.Get("/:id", h.GetByID)
	attributesRouter.Put("/:id", common.JwtMiddleware, h.Update)
	attributesRouter.Delete("/:id", common.JwtMiddleware, h.Delete)
}

// Attribute godoc
// @Summary Get attributes
// @Description Get all attributes
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {array} dtos.AttributePaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Router /attributes [get]
func (h *AttributeHandler) Fetch(c *fiber.Ctx) error {
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

	if attributes, err := h.service.Fetch(c.Context(), page, size, sortBy, orderBy); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(attributes)
	}
}

// Attribute godoc
// @Summary Get attribute by id
// @Description Get attribute by id
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} dtos.AttributeDto
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Router /attributes/{id} [get]
func (h *AttributeHandler) GetByID(c *fiber.Ctx) error {
	if id, err := c.ParamsInt("id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	} else {
		if attribute, err := h.service.GetByID(c.Context(), id); err != nil {
			switch err {
			case common.ErrNotFound:
				return c.SendStatus(fiber.StatusNotFound)
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(err)
			}
		} else {
			return c.JSON(attribute)
		}
	}
}

// Attribute godoc
// @Summary Update attribute
// @Description Update attribute by id
// @Tags attributes
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param dto body dtos.UpdateAttributeDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /attributes/{id} [put]
func (h *AttributeHandler) Update(c *fiber.Ctx) error {
	var body dtos.UpdateAttributeDto
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

// Attribute godoc
// @Summary Create attribute
// @Description Create new attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.CreateAttributeDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /attributes [post]
func (h *AttributeHandler) Create(c *fiber.Ctx) error {
	var body dtos.CreateAttributeDto
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

// Attribute godoc
// @Summary Delete attribute
// @Description Delete attribute by id
// @Tags attributes
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param Authorization header string true "Bearer"
// @Router /attributes/{id} [delete]
func (h *AttributeHandler) Delete(c *fiber.Ctx) error {
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

// Attribute godoc
// @Summary Search attribute
// @Description Search attributes by attribute name
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {array} dtos.AttributePaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param q query string true "query string to search in name"
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Router /attributes/search [get]
func (h *AttributeHandler) Search(c *fiber.Ctx) error {
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
