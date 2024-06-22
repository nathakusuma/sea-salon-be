package main

import (
	"github.com/joho/godotenv"
	"github.com/nathakusuma/sea-salon-be/internal/app/config"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	env := os.Getenv("ENV")
	if err != nil && env == "" {
		log.Fatalln("fail to load env")
	}

	db := config.NewDatabase()
	app := config.NewFiber()

	config.StartApp(&config.StartAppConfig{
		DB:  db,
		App: app,
	})

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		log.Fatalln(err)
	}
}
