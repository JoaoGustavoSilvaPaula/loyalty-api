package main

import (
	"os"

	"github.com/joaogustavosp/loyalty-api/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão
	}

	app := app.NewApp()
	app.Run(":" + port)

}
