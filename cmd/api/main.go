package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/joho/godotenv"

	"mob-backend/internal/config"
	"mob-backend/internal/handler/auth"
	"mob-backend/internal/middleware"
	repo "mob-backend/internal/repository/auth"
	email "mob-backend/internal/service/email"
	usecase "mob-backend/internal/usecase/auth"

	docs "mob-backend/docs"
)

// @title Merchant Onboarding Auth API
// @version 1.0
// @description Auth endpoints untuk merchant onboarding (register, verify email, login, me)
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// main bootstraps the HTTP server and wires auth routes with clean architecture layers.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env not loaded:", err)
	}
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Default is only for local development to keep the app runnable.
		jwtSecret = "change-this-secret"
		log.Println("JWT_SECRET not set, using a default development secret")
	}

	dsn := "root:root@tcp(127.0.0.1:3306)/onboarding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := config.NewMySQLDB(dsn)
	if err != nil {
		log.Println("failed to connect database:", err)
		return
	}

	userRepo := repo.NewUserRepository(db)
	tokenRepo := repo.NewEmailVerificationRepository(db)

	postmarkToken := os.Getenv("POSTMARK_SERVER_TOKEN")
	postmarkFrom := os.Getenv("POSTMARK_FROM_EMAIL")
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	emailService := email.NewPostmarkEmailService(postmarkToken, postmarkFrom, appURL)
	if postmarkToken == "" || postmarkFrom == "" {
		log.Println("Postmark env vars incomplete; verification emails will log errors until configured")
	}

	authUsecase := usecase.NewAuthService(userRepo, tokenRepo, emailService, usecase.TokenConfig{
		Secret:         jwtSecret,
		AccessTokenTTL: time.Hour,
		Issuer:         "mob-backend",
	})
	authHandler := auth.NewHandler(authUsecase)

	authGroup := r.Group("/api/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.GET("/verify-email", authHandler.VerifyEmail)
	authGroup.POST("/login", authHandler.Login)
	authGroup.GET("/me", middleware.AuthMiddleware(jwtSecret), authHandler.Me)

	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
