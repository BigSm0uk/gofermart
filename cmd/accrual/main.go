package main

import (
	"log"

	"github.com/BigSm0uk/gofermart/internal/app/accrual"
)

func main() {
	app, err := accrual.InitApp()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
