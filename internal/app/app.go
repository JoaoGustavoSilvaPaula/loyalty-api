package app

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/joaogustavosp/loyalty-api/internal/app/handlers/v1"
	"github.com/joaogustavosp/loyalty-api/internal/app/middleware"
	"github.com/joaogustavosp/loyalty-api/internal/db/mongodb"
	"github.com/joaogustavosp/loyalty-api/internal/services"
)

type App struct {
	Router *gin.Engine
}

func NewApp() *App {
	db := mongodb.GetDatabase("loyalty")
	userCollection := db.Collection("users")

	userService := services.NewUserService(userCollection)
	extractService := services.NewExtractService()

	userHandler := v1.NewUserHandler(userService)
	authHandler := v1.NewAuthHandler(userService)
	extractHandler := v1.NewExtractHandler(extractService)

	router := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	router.Use(cors.New(config))

	v1Route := router.Group("/api/v1")
	{
		v1Route.POST("/users", userHandler.CreateUser)
		v1Route.GET("/users/:id", middleware.AuthMiddleware(), userHandler.GetUser)
		v1Route.PUT("/users/:id", middleware.AuthMiddleware(), userHandler.UpdateUser)
		v1Route.DELETE("/users/:id", middleware.AuthMiddleware(), userHandler.DeleteUser)
		v1Route.GET("/users/exists/:cpf", userHandler.UserExistsByCPF)
		v1Route.POST("/users/create-password", userHandler.CreatePassword)

		auth := v1Route.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.GET("/validate-token", authHandler.ValidateToken)
		}
		v1Route.GET("/extract", extractHandler.ExtractData)
	}

	return &App{
		Router: router,
	}
}

func (a *App) Run(addr string) {
	log.Fatal(a.Router.Run(addr))
}
