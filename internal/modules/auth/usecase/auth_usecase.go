package usecase

import (
	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/utils"
	"sunchem-backend/internal/modules/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo      domain.IUserRepository
	jwtSecret string
}

func NewAuthUseCase(repo domain.IUserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{repo: repo, jwtSecret: jwtSecret}
}

func (uc *AuthUseCase) Login(username, password string) (string, *domain.User, *errors.AppError) {
	user, err := uc.repo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.NewError(401, "INVALID_CREDENTIALS", "Tên đăng nhập hoặc mật khẩu không đúng")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.NewError(401, "INVALID_CREDENTIALS", "Tên đăng nhập hoặc mật khẩu không đúng")
	}
	token, err := utils.GenerateToken(uc.jwtSecret, user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, errors.Wrap(err, 500, "TOKEN_ERROR", "Lỗi tạo token")
	}
	return token, user, nil
}
