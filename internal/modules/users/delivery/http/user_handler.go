package http

import (
	"strconv"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/response"
	"sunchem-backend/internal/modules/auth/domain"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo domain.IUserRepository
}

func NewUserHandler(repo domain.IUserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.repo.FindAll()
	if err != nil {
		_ = c.Error(errors.Wrap(err, 500, "DB_ERROR", "Lỗi truy vấn"))
		return
	}
	response.Success(c, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if req.Role == "" {
		req.Role = "employee"
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		_ = c.Error(errors.Wrap(err, 500, "HASH_ERROR", "Lỗi mã hóa"))
		return
	}
	user := &domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     req.Role,
	}
	if err := h.repo.Create(user); err != nil {
		_ = c.Error(errors.Wrap(err, 500, "DB_ERROR", "Lỗi tạo tài khoản"))
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	user, err := h.repo.FindByID(uint(id))
	if err != nil {
		_ = c.Error(errors.Wrap(err, 404, "NOT_FOUND", "Không tìm thấy"))
		return
	}
	var req struct {
		Name     string `json:"name"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			_ = c.Error(errors.Wrap(err, 500, "HASH_ERROR", "Lỗi mã hóa mật khẩu"))
			return
		}
		user.Password = string(hashedPassword)
	}
	if err := h.repo.Update(user); err != nil {
		_ = c.Error(errors.Wrap(err, 500, "DB_ERROR", "Lỗi cập nhật"))
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}
	currentUserID := c.GetUint("userID")
	if uint(id) == currentUserID {
		_ = c.Error(errors.NewError(400, "SELF_DELETE", "Không thể xóa chính mình"))
		return
	}
	if err := h.repo.Delete(uint(id)); err != nil {
		_ = c.Error(errors.Wrap(err, 500, "DB_ERROR", "Lỗi xóa tài khoản"))
		return
	}
	response.SuccessWithMessage(c, "Đã xóa", nil)
}
