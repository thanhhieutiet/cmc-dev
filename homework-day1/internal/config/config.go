package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config chứa toàn bộ cấu hình ứng dụng
type Config struct {
	DB     DBConfig
	Server ServerConfig
}

// DBConfig chứa thông tin kết nối database
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// ServerConfig chứa cấu hình server
type ServerConfig struct {
	Port string
}

// DSN trả về connection string để kết nối PostgreSQL
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}

// Load đọc cấu hình từ file .env và environment variables
func Load() *Config {
	// Cố gắng đọc file .env, nếu không có thì dùng env vars hệ thống
	_ = godotenv.Load()

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "mini_asm"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

// getEnv đọc biến môi trường, trả về giá trị mặc định nếu không tìm thấy
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
