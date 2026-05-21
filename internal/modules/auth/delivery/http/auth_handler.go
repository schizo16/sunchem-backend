package http

import (
	"fmt"
	"net/http"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
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

// GetConfig returns the OIDC configuration for the frontend
func (h *AuthHandler) GetConfig(c *gin.Context) {
	cfg := h.uc.GetOIDCConfig()
	response.Success(c, cfg)
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

// Token is an alias for GenoractCallback — exchanges OIDC code for a full token bundle
func (h *AuthHandler) Token(c *gin.Context) {
	var req genoractCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	bundle, user, appErr := h.uc.Token(req.Code, req.RedirectURI)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	_ = user
	response.Success(c, bundle)
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken validates a refresh token and issues a new token bundle
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	bundle, appErr := h.uc.RefreshToken(req.RefreshToken)
	if appErr != nil {
		_ = c.Error(appErr)
		return
	}
	response.Success(c, bundle)
}

// UsersMe returns the current user profile in the format expected by the blog-admin
func (h *AuthHandler) UsersMe(c *gin.Context) {
	userID := c.GetUint("userID")
	username := c.GetString("username")
	role := c.GetString("role")
	response.Success(c, gin.H{
		"id":          fmt.Sprintf("%d", userID),
		"username":    username,
		"name":        username,
		"role":        role,
		"permissions": []string{"*"},
	})
}


