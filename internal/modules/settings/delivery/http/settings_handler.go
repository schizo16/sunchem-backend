package http

import (
	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/settings/usecase"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	uc *usecase.SettingUseCase
}

func NewSettingHandler(uc *usecase.SettingUseCase) *SettingHandler {
	return &SettingHandler{uc: uc}
}

func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, appErr := h.uc.GetAll()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, settings)
}

func (h *SettingHandler) Save(c *gin.Context) {
	var values map[string]string
	if err := c.ShouldBindJSON(&values); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if appErr := h.uc.Save(values); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.SuccessWithMessage(c, "Đã lưu", nil)
}
