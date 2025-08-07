package main

import (
	"event-processing-pipeline/internal/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	ginRouter := config.Engine()
	ginRouter = config.Routers(ginRouter)

	err := ginRouter.Run(":9000")

	if err != nil {
		log.Fatal(err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
