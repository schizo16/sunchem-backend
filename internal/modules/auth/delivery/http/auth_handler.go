package http

import (
	"net/http"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/modules/auth/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	token, user, appErr := h.uc.Login(req.Username, req.Password)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id":  c.GetUint("userID"),
		"username": c.GetString("username"),
		"role":     c.GetString("role"),
	})
}

type genoractCallbackReq struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri" binding:"required"`
}

func (h *AuthHandler) GenoractCallback(c *gin.Context) {
	var req genoractCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	token, user, appErr := h.uc.GenoractCallback(req.Code, req.RedirectURI)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
