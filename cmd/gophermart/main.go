package main

import (
	"log"

	"github.com/BigSm0uk/gofermart/internal/app/gophermart"
)

func main() {
	app, err := gophermart.InitApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	app.Run()
}
