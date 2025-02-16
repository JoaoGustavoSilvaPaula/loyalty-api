package main

import "github.com/joaogustavosp/loyalty-api/internal/app"

func main() {
	app := app.NewApp()
	app.Run(":8080")
}
