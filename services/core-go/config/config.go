package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv         string
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	RedisHost      string
	RedisPort      string
	RedisPassword  string
	RabbitMQURL    string
	AIServiceURL   string
	JWTSecret      string
	JWTExpireHours int
}

func LoadConfig() *Config {
	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "9000"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "Cebes1234"),
		DBName:         getEnv("DB_NAME", "bank_sampah_db"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", "redis123"),
		RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		AIServiceURL:   getEnv("AI_SERVICE_URL", "http://localhost:8000"),
		JWTSecret:      getEnv("JWT_SECRET", "supersecretbankkeys1234567890abcdef"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
