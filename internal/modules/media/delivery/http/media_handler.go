package http

import (
	"strconv"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/media/usecase"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	uc *usecase.MediaUseCase
}

func NewMediaHandler(uc *usecase.MediaUseCase) *MediaHandler {
	return &MediaHandler{uc: uc}
}

func (h *MediaHandler) List(c *gin.Context) {
	files, appErr := h.uc.List()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, files)
}

func (h *MediaHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		_ = c.Error(errors.NewError(400, "NO_FILE", "Vui lòng chọn file"))
		return
	}
	defer file.Close()

	media, appErr := h.uc.Upload(header.Filename, file, header.Size)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, media)
}

func (h *MediaHandler) Delete(c *gin.Context) {
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
