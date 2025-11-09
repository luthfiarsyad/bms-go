package util

import (
	"fmt"
	"log"
	"os"
	"testing"

	"bms-go/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitTestDB() *gorm.DB {
	host := os.Getenv("DB_HOST_TEST")
	user := os.Getenv("DB_USER_TEST")
	pass := os.Getenv("DB_PASS_TEST")
	name := os.Getenv("DB_NAME_TEST")

	if host == "" || user == "" || name == "" {
		log.Fatal("Test database environment variables not set. Please define DB_HOST_TEST, DB_USER_TEST, DB_PASS_TEST, DB_NAME_TEST")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect test database: %v", err)
	}

	db.Exec("DROP TABLE IF EXISTS favorites, books;")
	db.AutoMigrate(&model.Book{}, &model.Favorite{})
	return db
}

func TestInitTestDB(t *testing.T) {
	db := InitTestDB()
	if db == nil {
		t.Fatal("expected db connection, got nil")
	}
}
