package db

import (
	"fmt"
	"log"

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

func InitDatabase() error {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Name)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PoolSize: 10, // adjust the pool size as needed
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to database")

	return nil
}

func GetDB() *gorm.DB {
	return db
}