package main

import (
	"log"

	"distrise/internal/app"
)

func main() {
	server, err := app.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
