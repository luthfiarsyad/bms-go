package main

import (
	"bms-go/config"
	"bms-go/docs"
	"bms-go/internal/infra/handler"
	"bms-go/internal/infra/repository"
	"bms-go/internal/service"
	"bms-go/util"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Book Management System API
// @version 1.0
// @description REST API sederhana untuk mengelola buku dan daftar favorit pengguna.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bms-go.local

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	config.LoadEnv()

	db := util.InitDB()

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	favRepo := repository.NewFavoriteRepository(db)
	favService := service.NewFavoriteService(favRepo, bookRepo)
	favHandler := handler.NewFavoriteHandler(favService)

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	bookHandler.RegisterRoutes(r)
	favHandler.RegisterRoutes(r)

	r.NoRoute(handler.NotFoundHandler)

	log.Println("Server running at http://localhost:8080")
	log.Println("Swagger docs available at http://localhost:8080/swagger/index.html")

	// Run server
	r.Run(":8080")
}
