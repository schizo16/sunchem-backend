package http

import (
	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/settings/usecase"

	"github.com/gin-gonic/gin"
)

// SEOHandler reuses the settings use-case, storing SEO values under "seo." prefixed keys.
type SEOHandler struct {
	uc *usecase.SettingUseCase
}

func NewSEOHandler(uc *usecase.SettingUseCase) *SEOHandler {
	return &SEOHandler{uc: uc}
}

// GetSEO returns all settings whose keys start with "seo." (stripped of the prefix)
func (h *SEOHandler) GetSEO(c *gin.Context) {
	all, appErr := h.uc.GetAll()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	result := make(map[string]string)
	const prefix = "seo."
	for k, v := range all {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			result[k[len(prefix):]] = v
		}
	}
	response.Success(c, result)
}

// SaveSEO persists SEO settings under "seo." prefixed keys
func (h *SEOHandler) SaveSEO(c *gin.Context) {
	var values map[string]string
	if err := c.ShouldBindJSON(&values); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	prefixed := make(map[string]string)
	for k, v := range values {
		prefixed["seo."+k] = v
	}
	if appErr := h.uc.Save(prefixed); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.SuccessWithMessage(c, "Đã lưu SEO", nil)
}
