package http

import (
	"os"

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

// GetGeneral returns general site settings (public endpoint)
func (h *SettingHandler) GetGeneral(c *gin.Context) {
	response.Success(c, gin.H{
		"site_name":    "Sunchem",
		"is_installed": true,
	})
}

// GetOIDC returns OIDC settings from env
func (h *SettingHandler) GetOIDC(c *gin.Context) {
	response.Success(c, gin.H{
		"authority":    os.Getenv("OIDC_AUTHORITY"),
		"client_id":    os.Getenv("OIDC_CLIENT_ID"),
		"redirect_uri": os.Getenv("OIDC_REDIRECT_URI"),
	})
}

// SaveOIDC stores OIDC settings (persisted as settings keys)
func (h *SettingHandler) SaveOIDC(c *gin.Context) {
	var values map[string]string
	if err := c.ShouldBindJSON(&values); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	// prefix keys with "oidc."
	prefixed := make(map[string]string)
	for k, v := range values {
		prefixed["oidc."+k] = v
	}
	if appErr := h.uc.Save(prefixed); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.SuccessWithMessage(c, "Đã lưu OIDC", nil)
}

// GetStorage returns storage settings
func (h *SettingHandler) GetStorage(c *gin.Context) {
	settings, appErr := h.uc.GetAll()
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	// filter keys prefixed with "storage."
	result := make(map[string]string)
	for k, v := range settings {
		if len(k) > 8 && k[:8] == "storage." {
			result[k[8:]] = v
		}
	}
	response.Success(c, result)
}

// SaveStorage stores storage settings
func (h *SettingHandler) SaveStorage(c *gin.Context) {
	var values map[string]string
	if err := c.ShouldBindJSON(&values); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	prefixed := make(map[string]string)
	for k, v := range values {
		prefixed["storage."+k] = v
	}
	if appErr := h.uc.Save(prefixed); appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.SuccessWithMessage(c, "Đã lưu storage", nil)
}
