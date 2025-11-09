package util

import (
	"bms-go/internal/model"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// requiredKeys berisi daftar key yang wajib diisi
var requiredKeys = []string{
	"database.user",
	"database.host",
	"database.port",
	"database.name",
}

func InitDB() *gorm.DB {
	// Setup Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, attempting to load from environment variables: %v", err)
	}

	viper.AutomaticEnv()

	missingKeys := []string{}
	for _, key := range requiredKeys {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		log.Fatalf("Missing required configuration values: %v", missingKeys)
	}

	user := viper.GetString("database.user")
	pass := viper.GetString("database.pass")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	name := viper.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	if err := db.AutoMigrate(&model.Book{}, &model.Favorite{}); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	log.Printf("Connected to MySQL [%s:%s] successfully!", host, name)
	return db
}
