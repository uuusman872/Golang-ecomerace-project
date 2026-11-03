package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
	SMTP     SMTPConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret                string
	ExpiresIn             time.Duration
	RefreshTokenExpiresIn time.Duration
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
	S3Endpoint      string
	EventQueueName  string
}

type UploadConfig struct {
	Path           string
	MaxUploadSize  int64
	UploadProvider string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	jwtSecret := getEnv("JWT_SECRET", "your_jwt_secret_key")
	jwtExpiresInStr := getEnv("JWT_EXPIRES_IN", "24h")
	refreshTokenExpiresInStr := getEnv("REFRESH_TOKEN_EXPIRES_IN", "72h")
	maxUploadSizeStr := getEnv("MAX_UPLOAD_SIZE", "10485760")
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "1025"))

	jwtExpiresIn, _ := time.ParseDuration(jwtExpiresInStr)
	refreshTokenExpires, _ := time.ParseDuration(refreshTokenExpiresInStr)
	maxUploadSize, _ := strconv.ParseInt(maxUploadSizeStr, 10, 64)

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "admin1234"),
			Name:     getEnv("DB_NAME", "ecomerace_shop"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:                jwtSecret,
			ExpiresIn:             jwtExpiresIn,
			RefreshTokenExpiresIn: refreshTokenExpires,
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:        getEnv("AWS_S3_BUCKET", "ecommerce-uploads"),
			S3Endpoint:      getEnv("AWS_S3_ENDPOINT", "http://localhost:4566"),
			EventQueueName:  getEnv("AWS_EVENT_QUEUE_NAME", "ecommerce-events"),
		},
		Upload: UploadConfig{
			Path:           getEnv("UPLOAD_PATH", "./uploads"),
			MaxUploadSize:  maxUploadSize,
			UploadProvider: getEnv("UPLOAD_PROVIDER", "local"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     smtpPort,
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@shop.com"),
		},
	}, nil
}
