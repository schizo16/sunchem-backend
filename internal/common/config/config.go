package config

import "os"

type Config struct {
	AppEnv     string
	DBType     string
	DBDSN      string
	ServerPort           string
	JWTSecret            string
	UploadDir            string
	GenoractClientID     string
	GenoractClientSecret string
	OIDCAuthority        string
	OIDCClientID         string
	OIDCRedirectURI      string
}

func LoadConfig() *Config {
	dsn := getEnv("DB_DSN", "sunchem.db")
	if dsn == "sunchem.db" {
		wd, _ := os.Getwd()
		dsn = wd + "/sunchem.db"
	}
	return &Config{
		AppEnv:     getEnv("APP_ENV", "dev"),
		DBType:     getEnv("DB_TYPE", "sqlite"),
		DBDSN:      dsn,
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		JWTSecret:            getEnv("JWT_SECRET", "sunchem-secret-key-change-in-production"),
		UploadDir:            getEnv("UPLOAD_DIR", "./uploads"),
		GenoractClientID:     getEnv("GENORACT_CLIENT_ID", ""),
		GenoractClientSecret: getEnv("GENORACT_CLIENT_SECRET", ""),
		OIDCAuthority:        getEnv("OIDC_AUTHORITY", ""),
		OIDCClientID:         getEnv("OIDC_CLIENT_ID", ""),
		OIDCRedirectURI:      getEnv("OIDC_REDIRECT_URI", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
