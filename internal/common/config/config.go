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
}

func LoadConfig() *Config {
	dsn := getEnv("DB_DSN", "sunchem.db")
	if dsn == "sunchem.db" {
		wd, _ := os.Getwd()
		dsn = wd + "/sunchem.db"
	}
	port := getEnv("PORT", "")
	if port == "" {
		port = getEnv("SERVER_PORT", "8080")
	}
	return &Config{
		AppEnv:     getEnv("APP_ENV", "dev"),
		DBType:     getEnv("DB_TYPE", "sqlite"),
		DBDSN:      dsn,
		ServerPort:           port,
		JWTSecret:            getEnv("JWT_SECRET", "sunchem-secret-key-change-in-production"),
		UploadDir:            getEnv("UPLOAD_DIR", "./uploads"),
		GenoractClientID:     getEnv("GENORACT_CLIENT_ID", ""),
		GenoractClientSecret: getEnv("GENORACT_CLIENT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
