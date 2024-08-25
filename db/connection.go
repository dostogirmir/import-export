// db/connection.go
package db

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_DATABASE"`
}

var (
	db  *gorm.DB
	err error
)

// InitDatabase initializes the database connection using environment variables
func InitDatabase() error {
	// Explicitly set the .env file location
	viper.SetConfigFile(".env")

	// Load the .env file
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Log the configuration to ensure values are correctly loaded
	log.Printf("DB Config: %+v\n", config)

	// Check if the necessary configuration is populated
	if config.Username == "" || config.Host == "" || config.Port == 0 || config.Name == "" {
		return fmt.Errorf("database configuration is incomplete")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Name)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from GORM DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(100)                // Increase to handle more concurrent connections
	sqlDB.SetMaxIdleConns(20)                 // Increase to keep more idle connections open
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Reduce lifetime to refresh stale connections

	log.Println("Connected to the database successfully")

	return nil
}

// GetDB returns the initialized GORM DB instance
func GetDB() *gorm.DB {
	return db
}