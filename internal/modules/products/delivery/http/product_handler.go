package http

import (
	"strconv"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/products/domain"
	"sunchem-backend/internal/modules/products/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	uc *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) List(c *gin.Context) {
	products, appErr := h.uc.List()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	product, appErr := h.uc.GetByID(uint(id))
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, product)
}

func (h *ProductHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	product, appErr := h.uc.GetBySlug(slug)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Create(&product); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	updated, appErr := h.uc.UpdatePartial(uint(id), updates)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, updated)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Delete(uint(id)); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.SuccessWithMessage(c, "Đã xóa", nil)
}
