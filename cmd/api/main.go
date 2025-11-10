package main

import (
	"file-sharing/internal/config"
	"file-sharing/internal/storage"
	"file-sharing/internal/transport/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// ---- 1. Khởi tạo Database Connection ----
	config, err := config.LoadConfig("./env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	log.Println("Connected to database")

	// ---- 2. Khởi tạo UserRepository ----
	userRepo := storage.NewUserRepository(db)

	// ---- 3. Khởi tạo Gin Router ----
	router := gin.Default()

	// ---- 4. Khởi tạo Handlers & Middlewares ----
	userHandler := http.NewUserHandler()
	authMiddleware := http.AuthMiddleware(userRepo)
	fileHandler := http.InitFileUploadHandler(db.DB)

	// ---- 5. Đăng ký API Routes ----
	api := router.Group("/api")
	{
		authed := api.Group("/")
		authed.Use(authMiddleware)
		{
			authed.GET("/me", userHandler.GetCurrentUser)
			authed.POST("/v1/files", fileHandler)
		}
	}

	// ---- Swagger UI ----
	router.StaticFile("/openapi.yaml", "./api/openapi.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.yaml")))

	// ---- 6. Khởi động HTTP Server ----
	log.Printf("Starting server on %s", config.HTTPServerAddress)
	if err := router.Run(config.HTTPServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
