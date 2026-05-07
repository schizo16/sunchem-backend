package http

import (
	"strconv"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/blog/domain"
	"sunchem-backend/internal/modules/blog/usecase"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	uc *usecase.BlogUseCase
}

func NewBlogHandler(uc *usecase.BlogUseCase) *BlogHandler {
	return &BlogHandler{uc: uc}
}

func (h *BlogHandler) List(c *gin.Context) {
	posts, appErr := h.uc.List()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, posts)
}

func (h *BlogHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	post, appErr := h.uc.GetByID(uint(id))
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, post)
}

func (h *BlogHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	post, appErr := h.uc.GetBySlug(slug)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, post)
}

func (h *BlogHandler) Create(c *gin.Context) {
	var post domain.BlogPost
	if err := c.ShouldBindJSON(&post); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if post.Status == "" {
		post.Status = "draft"
	}
	if appErr := h.uc.Create(&post); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, post)
}

func (h *BlogHandler) Update(c *gin.Context) {
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

func (h *BlogHandler) Delete(c *gin.Context) {
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
