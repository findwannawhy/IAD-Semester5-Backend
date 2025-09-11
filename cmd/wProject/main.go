package main

import (
	"log"

	"github.com/findwannawhy/IAD-Semester5/internal/api"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}