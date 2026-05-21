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
	"time"

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
	oidcAuthority        string
	oidcClientID         string
	oidcRedirectURI      string
}

func NewAuthUseCase(repo domain.IUserRepository, jwtSecret, clientID, clientSecret string) *AuthUseCase {
	return &AuthUseCase{
		repo:                 repo,
		jwtSecret:            jwtSecret,
		genoractClientID:     clientID,
		genoractClientSecret: clientSecret,
	}
}

// SetOIDCConfig sets the OIDC configuration (authority, clientID, redirectURI)
func (uc *AuthUseCase) SetOIDCConfig(authority, clientID, redirectURI string) {
	uc.oidcAuthority = authority
	uc.oidcClientID = clientID
	uc.oidcRedirectURI = redirectURI
}

// OIDCConfig holds the OIDC configuration returned to the frontend
type OIDCConfig struct {
	Authority   string `json:"authority"`
	ClientID    string `json:"client_id"`
	OIDCIssuer  string `json:"oidc_issuer"`
	OIDCClientID string `json:"oidc_client_id"`
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
}

// TokenBundle holds the full token response expected by the blog-admin
type TokenBundle struct {
	AccessToken    string `json:"access_token"`
	IDToken        string `json:"id_token"`
	IdentifyToken  string `json:"identify_token"`
	RefreshToken   string `json:"refresh_token"`
	Expiry         string `json:"expiry"`
}

// GetOIDCConfig returns the OIDC configuration
func (uc *AuthUseCase) GetOIDCConfig() *OIDCConfig {
	return &OIDCConfig{
		Authority:    uc.oidcAuthority,
		ClientID:     uc.oidcClientID,
		OIDCIssuer:   uc.oidcAuthority,
		OIDCClientID: uc.oidcClientID,
		RedirectURI:  uc.oidcRedirectURI,
		Scope:        "openid profile email offline_access",
	}
}

// buildTokenBundle creates the full token bundle given a user
func (uc *AuthUseCase) buildTokenBundle(user *domain.User) (*TokenBundle, error) {
	accessToken, err := utils.GenerateToken(uc.jwtSecret, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}
	permissions := []string{"*"}
	idToken, err := utils.GenerateIDToken(uc.jwtSecret, user.Username, permissions)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(uc.jwtSecret, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}
	expiry := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	return &TokenBundle{
		AccessToken:   accessToken,
		IDToken:       idToken,
		IdentifyToken: idToken,
		RefreshToken:  refreshToken,
		Expiry:        expiry,
	}, nil
}

// Token is an alias for GenoractCallback — exchanges OIDC code for tokens
func (uc *AuthUseCase) Token(code string, redirectURI string) (*TokenBundle, *domain.User, *errors.AppError) {
	_, user, appErr := uc.GenoractCallback(code, redirectURI)
	if appErr != nil {
		return nil, nil, appErr
	}
	bundle, err := uc.buildTokenBundle(user)
	if err != nil {
		return nil, nil, errors.Wrap(err, 500, "TOKEN_ERROR", "Lỗi tạo token bundle")
	}
	return bundle, user, nil
}

// RefreshToken validates a refresh token and issues a new token bundle
func (uc *AuthUseCase) RefreshToken(refreshTokenStr string) (*TokenBundle, *errors.AppError) {
	claims, err := utils.ParseToken(uc.jwtSecret, refreshTokenStr)
	if err != nil {
		return nil, errors.NewError(401, "INVALID_TOKEN", "Refresh token không hợp lệ")
	}
	user, err := uc.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.NewError(401, "INVALID_TOKEN", "Người dùng không tồn tại")
	}
	bundle, err := uc.buildTokenBundle(user)
	if err != nil {
		return nil, errors.Wrap(err, 500, "TOKEN_ERROR", "Lỗi tạo token bundle")
	}
	return bundle, nil
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

// LoginBundle authenticates with username/password and returns a full token bundle
func (uc *AuthUseCase) LoginBundle(username, password string) (*TokenBundle, *domain.User, *errors.AppError) {
	user, err := uc.repo.FindByUsername(username)
	if err != nil {
		return nil, nil, errors.NewError(401, "INVALID_CREDENTIALS", "Tên đăng nhập hoặc mật khẩu không đúng")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.NewError(401, "INVALID_CREDENTIALS", "Tên đăng nhập hoặc mật khẩu không đúng")
	}
	bundle, err := uc.buildTokenBundle(user)
	if err != nil {
		return nil, nil, errors.Wrap(err, 500, "TOKEN_ERROR", "Lỗi tạo token")
	}
	return bundle, user, nil
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
			Role:       "admin",
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
