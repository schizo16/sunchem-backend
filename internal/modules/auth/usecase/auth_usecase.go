package usecase

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"sunchem-backend/internal/common/errors"
	"sunchem-backend/internal/common/utils"
	"sunchem-backend/internal/modules/auth/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo                 domain.IUserRepository
	jwtSecret            string
	genoractClientID     string
	genoractClientSecret string
}

func NewAuthUseCase(repo domain.IUserRepository, jwtSecret, clientID, clientSecret string) *AuthUseCase {
	return &AuthUseCase{
		repo:                 repo,
		jwtSecret:            jwtSecret,
		genoractClientID:     clientID,
		genoractClientSecret: clientSecret,
	}
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

type genoractTokenResponse struct {
	IdToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
}

type jwksResponse struct {
	Keys []struct {
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		Crv string `json:"crv"`
		X   string `json:"x"`
	} `json:"keys"`
}

func fetchGenoractPublicKey(kid string) (ed25519.PublicKey, error) {
	resp, err := http.Get("https://accounts.genoract.com/api/v1/certs")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch certs: %v", err)
	}
	defer resp.Body.Close()

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode jwks: %v", err)
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid && key.Crv == "Ed25519" {
			xBytes, err := base64.RawURLEncoding.DecodeString(key.X)
			if err != nil {
				return nil, fmt.Errorf("failed to decode x: %v", err)
			}
			return ed25519.PublicKey(xBytes), nil
		}
	}
	return nil, fmt.Errorf("kid %s not found", kid)
}

func (uc *AuthUseCase) GenoractCallback(code string, redirectURI string) (string, *domain.User, *errors.AppError) {
	// 1. Exchange code for token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", uc.genoractClientID)
	data.Set("client_secret", uc.genoractClientSecret)

	req, _ := http.NewRequest("POST", "https://accounts.genoract.com/api/v1/oauth2/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, errors.Wrap(err, 500, "GENORACT_ERROR", "Lỗi kết nối tới Genoract")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", nil, errors.NewError(500, "GENORACT_ERROR", fmt.Sprintf("Lỗi từ Genoract: %s", string(bodyBytes)))
	}

	var tokenResp genoractTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", nil, errors.Wrap(err, 500, "GENORACT_ERROR", "Lỗi đọc dữ liệu từ Genoract")
	}

	// 2. Parse and verify id_token
	token, err := jwt.Parse(tokenResp.IdToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid missing in header")
		}
		return fetchGenoractPublicKey(kid)
	})

	if err != nil || !token.Valid {
		return "", nil, errors.Wrap(err, 401, "INVALID_TOKEN", "ID Token không hợp lệ")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, errors.NewError(401, "INVALID_TOKEN", "Claims không hợp lệ")
	}

	// 3. Verify iss and aud
	iss, _ := claims["iss"].(string)
	if iss != "https://accounts.genoract.com" {
		return "", nil, errors.NewError(401, "INVALID_TOKEN", "Issuer không hợp lệ")
	}

	audValid := false
	if audStr, ok := claims["aud"].(string); ok {
		if audStr == uc.genoractClientID {
			audValid = true
		}
	} else if audArr, ok := claims["aud"].([]interface{}); ok {
		for _, a := range audArr {
			if aStr, okStr := a.(string); okStr && aStr == uc.genoractClientID {
				audValid = true
				break
			}
		}
	}
	
	if !audValid {
		return "", nil, errors.NewError(401, "INVALID_TOKEN", "Audience không hợp lệ")
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", nil, errors.NewError(401, "INVALID_TOKEN", "Subject bị thiếu")
	}

	// 4. Find or Create User
	user, err := uc.repo.FindByGenoractID(sub)
	if err != nil {
		user = &domain.User{
			GenoractID: &sub,
			Username:   "genoract_" + sub,
			Name:       "Genoract User",
			Role:       "employee",
		}
		if err := uc.repo.Create(user); err != nil {
			return "", nil, errors.Wrap(err, 500, "DB_ERROR", "Lỗi tạo người dùng")
		}
	}

	// 5. Generate System JWT
	sysToken, err := utils.GenerateToken(uc.jwtSecret, user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, errors.Wrap(err, 500, "TOKEN_ERROR", "Lỗi tạo token hệ thống")
	}

	return sysToken, user, nil
}
