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
	userHandler := v1.NewUserHandler(userService)
	authHandler := v1.NewAuthHandler(userService)

	router := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	router.Use(cors.New(config))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users/:id", middleware.AuthMiddleware(), userHandler.GetUser)
		v1.PUT("/users/:id", middleware.AuthMiddleware(), userHandler.UpdateUser)
		v1.DELETE("/users/:id", middleware.AuthMiddleware(), userHandler.DeleteUser)
		v1.GET("/users/exists/:cpf", userHandler.UserExistsByCPF)
		v1.POST("/users/create-password", userHandler.CreatePassword)

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.GET("/validate-token", authHandler.ValidateToken)
		}
	}

	return &App{
		Router: router,
	}
}

func (a *App) Run(addr string) {
	log.Fatal(a.Router.Run(addr))
}
