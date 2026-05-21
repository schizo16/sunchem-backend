package http

import (
	"strconv"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/categories/domain"
	"sunchem-backend/internal/modules/categories/usecase"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	uc *usecase.CategoryUseCase
}

func NewCategoryHandler(uc *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) List(c *gin.Context) {
	cats, appErr := h.uc.List()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, cats)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var cat domain.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Create(&cat); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	var cat domain.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Update(uint(id), &cat); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, cat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
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
