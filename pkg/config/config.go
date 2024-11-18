package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
}

type ServerConfig struct {
	Port            string
	Host            string
	ShutdownTimeout time.Duration // Time to wait for graceful shutdown
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConnections  int           // Maximum number of database connections
	ConnTimeout     time.Duration // Connection timeout
	MaxIdleConns    int           // Maximum number of idle connections
	MaxOpenConns    int           // Maximum number of open connections
	ConnMaxLifetime time.Duration // Maximum lifetime of connections
}

type APIConfig struct {
	DefaultPageSize    int
	MaxPageSize       int
	RateLimitPerMin   int
	RequestTimeout    time.Duration // Timeout for API requests
	ReadTimeout       time.Duration // Server read timeout
	WriteTimeout      time.Duration // Server write timeout
	ShutdownTimeout   time.Duration // Graceful shutdown timeout
	EnableSwagger     bool          // Enable Swagger documentation
	EnablePrometheus  bool          // Enable Prometheus metrics
	EnableHealthCheck bool          // Enable health check endpoint
}

// LoadConfig returns a new Config struct populated with values from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Host:            getEnv("SERVER_HOST", "localhost"),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", "5s"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "postgres"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "user_api"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxConnections:  getEnvAsInt("DB_MAX_CONNECTIONS", 100),
			ConnTimeout:     getEnvAsDuration("DB_CONN_TIMEOUT", "10s"),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "1h"),
		},
		API: APIConfig{
			DefaultPageSize:   getEnvAsInt("API_DEFAULT_PAGE_SIZE", 10),
			MaxPageSize:      getEnvAsInt("API_MAX_PAGE_SIZE", 100),
			RateLimitPerMin:  getEnvAsInt("API_RATE_LIMIT_PER_MIN", 60),
			RequestTimeout:   getEnvAsDuration("API_REQUEST_TIMEOUT", "30s"),
			ReadTimeout:      getEnvAsDuration("API_READ_TIMEOUT", "15s"),
			WriteTimeout:     getEnvAsDuration("API_WRITE_TIMEOUT", "15s"),
			ShutdownTimeout:  getEnvAsDuration("API_SHUTDOWN_TIMEOUT", "5s"),
			EnableSwagger:    getEnvAsBool("API_ENABLE_SWAGGER", true),
			EnablePrometheus: getEnvAsBool("API_ENABLE_PROMETHEUS", true),
			EnableHealthCheck: getEnvAsBool("API_ENABLE_HEALTH_CHECK", true),
		},
	}
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return splitAndTrim(value, ",")
	}
	return defaultValue
}

// Helper function to split and trim a string
func splitAndTrim(s string, sep string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, sep)
	var result []string
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
