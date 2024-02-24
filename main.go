package main

import (
	"os"

	"github.com/joho/godotenv"

	"gpu/app"
)

func main() {
	godotenv.Load()

	a := app.App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("JWT_SECRET"),
		os.Getenv("STRIPE_SECRET"),
		os.Getenv("STRIPE_WEBHOOK"),
		os.Getenv("DEV") == "true")

	a.Run(":8080")
}
