package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig creates a new database configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "user_api"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect establishes a connection to the database with retries
func Connect(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	var db *sql.DB
	var err error
	
	// Try to connect with retries
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		}

		// Test the connection
		err = db.Ping()
		if err == nil {
			return db, nil
		}

		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(time.Second * 2)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
