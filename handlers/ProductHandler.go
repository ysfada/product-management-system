package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ysfada/product-management-system/domain/common"
	"github.com/ysfada/product-management-system/domain/dtos"
	"github.com/ysfada/product-management-system/domain/interfaces"
)

type ProductHandler struct {
	service interfaces.IProductService
}

func NewProductHandler(service interfaces.IProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

var _ interfaces.IProductHandler = (*ProductHandler)(nil)

func (h *ProductHandler) UseHandler(r fiber.Router) {
	productsRouter := r.Group("products")

	productsRouter.Get("/", h.Fetch)
	productsRouter.Post("/", h.Create)
	productsRouter.Get("/search", h.Search)
	productsRouter.Get("/:id", h.GetByID)
	productsRouter.Put("/:id", h.Update)
	productsRouter.Delete("/:id", common.JwtMiddleware, h.Delete)
	productsRouter.Get("/:id/images", h.GetImages)
	productsRouter.Post("/:id/images", common.JwtMiddleware, h.AddImage)
	productsRouter.Delete("/:id/images/:imageID", common.JwtMiddleware, h.RemoveImage)
	productsRouter.Get("/:id/variants", h.FetchVariants)
	productsRouter.Post("/:id/variants", common.JwtMiddleware, h.CreateVariant)
	productsRouter.Get("/:id/variants/search", h.SearchVariants)
	productsRouter.Get("/:id/variants/:variantID", h.GetVariantByID)
	productsRouter.Put("/:id/variants/:variantID", common.JwtMiddleware, h.UpdateVariant)
	productsRouter.Delete("/:id/variants/:variantID", common.JwtMiddleware, h.DeleteVariant)
	productsRouter.Get("/:id/variants/:variantID/attributes", h.GetAttributes)
	productsRouter.Post("/:id/variants/:variantID/attributes", common.JwtMiddleware, h.AddAttribute)
	productsRouter.Delete("/:id/variants/:variantID/attributes/:attributeID", common.JwtMiddleware, h.RemoveAttribute)
}

// Product godoc
// @Summary Get products
// @Description Get all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ProductPaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Router /products [get]
func (h *ProductHandler) Fetch(c *fiber.Ctx) error {
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

	if products, err := h.service.Fetch(c.Context(), page, size, sortBy, orderBy); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(products)
	}
}

// Product godoc
// @Summary Get product by id
// @Description Get product by id
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dtos.ProductDto
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	if id, err := c.ParamsInt("id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	} else {
		if product, err := h.service.GetByID(c.Context(), id); err != nil {
			switch err {
			case common.ErrNotFound:
				return c.SendStatus(fiber.StatusNotFound)
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(err)
			}
		} else {
			return c.JSON(product)
		}
	}
}

// Product godoc
// @Summary Update product
// @Description Update product by id
// @Tags products
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param dto body dtos.UpdateProductDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	var body dtos.UpdateProductDto
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

// Product godoc
// @Summary Create product
// @Description Create new product
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.CreateProductDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /products [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var body dtos.CreateProductDto
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
	return c.SendStatus(fiber.StatusCreated)
}

// Product godoc
// @Summary Delete product
// @Description Delete product by id
// @Tags products
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param Authorization header string true "Bearer"
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
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

// Product godoc
// @Summary Search product
// @Description Search products by product name
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ProductPaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param q query string true "query string to search in name"
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Router /products/search [get]
func (h *ProductHandler) Search(c *fiber.Ctx) error {
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

	if products, err := h.service.Search(c.Context(), q, page, size, sortBy, orderBy); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(products)
	}
}

// Product godoc
// @Summary Get images belongs product
// @Description Get images belongs product
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ImageDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Router /products/{id}/images [get]
func (h *ProductHandler) GetImages(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if images, err := h.service.GetImages(c.Context(), id); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(images)
	}
}

// Product godoc
// @Summary Add image to product
// @Description Add new image to product
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param image formData file true "product image"
// @Param id path int true "id"
// @Param Authorization header string true "Bearer"
// @Router /products/{id}/images [post]
func (h *ProductHandler) AddImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if file, err := c.FormFile("image"); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	} else {
		if err := h.service.AddImage(c.Context(), id, file); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		} else {
			return c.SendStatus(fiber.StatusCreated)
		}
	}
}

// Product godoc
// @Summary Remove image from product
// @Description Remove an image from product
// @Tags products
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param imageID path int true "imageID"
// @Param Authorization header string true "Bearer"
// @Router /products/{id}/images/{imageID} [delete]
func (h *ProductHandler) RemoveImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	imageID, err := c.ParamsInt("imageID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.service.RemoveImage(c.Context(), id, imageID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Product godoc
// @Summary Get product variants
// @Description Get product variants
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ProductVariantPaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "name or id"
// @Param orderBy query string false "ASC or DESC"
// @Param id path int true "id"
// @Router /products/{id}/variants [get]
func (h *ProductHandler) FetchVariants(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
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

	if products, err := h.service.FetchVariants(c.Context(), id, page, size, sortBy, orderBy); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(products)
	}
}

// Product godoc
// @Summary Get product variant by id
// @Description Get product variant by id
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dtos.ProductVariantDto
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param variantID path int true "variantID"
// @Router /products/{id}/variants/{variantID} [get]
func (h *ProductHandler) GetVariantByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	variantID, err := c.ParamsInt("variantID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if product, err := h.service.GetVariantByID(c.Context(), id, variantID); err != nil {
		switch err {
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(product)
	}
}

// Product godoc
// @Summary Create product variant
// @Description Create new product variant
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.CreateProductVariantDto true "dto"
// @Param id path int true "id"
// @Param Authorization header string true "Bearer"
// @Router /products/{id}/variants [post]
func (h *ProductHandler) CreateVariant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	var body dtos.CreateProductVariantDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if id != body.ProductId {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.CreateVariant(c.Context(), &body); err != nil {
		switch err {
		case common.ErrBadParamInput:
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusCreated)
}

// Product godoc
// @Summary Update product variant
// @Description Update product variant by id
// @Tags products
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param variantID path int true "variantID"
// @Param dto body dtos.UpdateProductVariantDto true "dto"
// @Param Authorization header string true "Bearer"
// @Router /products/{id}/variants/{variantID} [put]
func (h *ProductHandler) UpdateVariant(c *fiber.Ctx) error {
	var body dtos.UpdateProductVariantDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if c.Params("id") != fmt.Sprint(body.ProductId) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if c.Params("variantID") != fmt.Sprint(body.ID) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.UpdateVariant(c.Context(), &body); err != nil {
		switch err {
		case common.ErrBadParamInput:
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Product godoc
// @Summary Delete product variant
// @Description Delete product variant by id
// @Tags products
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param variantID path int true "variantID"
// @Param Authorization header string true "Bearer"
// @Router /products/{id}/variants/{variantID} [delete]
func (h *ProductHandler) DeleteVariant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	variantID, err := c.ParamsInt("variantID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.service.DeleteVariant(c.Context(), id, variantID); err != nil {
		switch err {
		case common.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Product godoc
// @Summary Get attributes belongs product variant
// @Description Get attributes belongs product variant
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} dtos.AttributeDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param variantID path int true "variantID"
// @Router /products/{id}/variants/{variantID}/attributes [get]
func (h *ProductHandler) GetAttributes(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	variantID, err := c.ParamsInt("variantID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if variants, err := h.service.GetAttributes(c.Context(), id, variantID); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(variants)
	}
}

// Product godoc
// @Summary Add attribute to product variant
// @Description Add new attribute to product variant
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param dto body dtos.CreateProductVariantAttributeDto true "dto"
// @Param id path int true "id"
// @Param variantID path int true "variantID"
// @Param Authorization header string true "Bearer"
// @Router /products/{id}/variants/{variantID}/attributes [post]
func (h *ProductHandler) AddAttribute(c *fiber.Ctx) error {
	// id, err := c.ParamsInt("id")
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(err)
	// }

	variantID, err := c.ParamsInt("variantID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	var body dtos.CreateProductVariantAttributeDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if variantID != body.ProductVariantID {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.AddAttribute(c.Context(), &body); err != nil {
		switch err {
		case common.ErrBadParamInput:
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}
	return c.SendStatus(fiber.StatusCreated)
}

// Product godoc
// @Summary Remove attribute from product variant
// @Description Remove an attribute from product variant
// @Tags products
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param id path int true "id"
// @Param variantID path int true "variantID"
// @Param attributeID path int true "attributeID"
// @Param Authorization header string true "Bearer"
// @Router /products/{id}/variants/{variantID}/attributes/{attributeID} [delete]
func (h *ProductHandler) RemoveAttribute(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	variantID, err := c.ParamsInt("variantID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	attributeID, err := c.ParamsInt("attributeID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.service.RemoveAttribute(c.Context(), id, variantID, attributeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Product godoc
// @Summary Search product variants
// @Description Search product variants
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ProductVariantPaginatedDto
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Param attrs query string false "ex: [{'type':'color', 'names': ['red', 'yellow']}, {'type':'size', 'names': ['36']}]"
// @Param q query string true "query string to search in name"
// @Param page query int false "page number"
// @Param size query int false "rows per page"
// @Param sortBy query string false "id, name, price"
// @Param orderBy query string false "ASC or DESC"
// @Param id path int true "id"
// @Router /products/{id}/variants/search [get]
func (h *ProductHandler) SearchVariants(c *fiber.Ctx) error {
	q := c.Query("q")
	if len(q) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	attrsStr := c.Query("attrs")
	var attrs []*dtos.AttributeSearchQueryDto
	if len(attrsStr) > 0 {
		if err := json.Unmarshal([]byte(attrsStr), &attrs); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
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
	if sortBy != "id" && sortBy != "name" && sortBy != "price" {
		sortBy = "id"
	}
	orderBy := strings.ToUpper(c.Query("orderBy", "ASC"))
	if orderBy != "ASC" && orderBy != "DESC" {
		orderBy = "ASC"
	}

	if products, err := h.service.SearchVariants(c.Context(), q, id, page, size, sortBy, orderBy, attrs); err != nil {
		switch err {
		// case common.ErrNotFound:
		// 	return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	} else {
		return c.JSON(products)
	}
}
