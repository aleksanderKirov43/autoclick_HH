package main

import (
	"log"

	"autoclick_HH/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
