package http

import (
	"strconv"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/tags/domain"
	"sunchem-backend/internal/modules/tags/usecase"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	uc *usecase.TagUseCase
}

func NewTagHandler(uc *usecase.TagUseCase) *TagHandler {
	return &TagHandler{uc: uc}
}

func (h *TagHandler) List(c *gin.Context) {
	tags, appErr := h.uc.List()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, tags)
}

func (h *TagHandler) Create(c *gin.Context) {
	var tag domain.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Create(&tag); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, tag)
}

func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	var tag domain.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Update(uint(id), &tag); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
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
