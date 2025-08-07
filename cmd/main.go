package main

import (
	"event-processing-pipeline/internal/config"
	"log"
)

func main() {
	ginRouter := config.Engine()
	ginRouter = config.Routers(ginRouter)

	err := ginRouter.Run(":9000")

	if err != nil {
		log.Fatal(err)
	}
}
