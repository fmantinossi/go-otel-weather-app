package main

import (
	"context"
	"log"
	"net/http"

	"service-b/handlers"
	"service-b/otel"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found, following with environment variables")
	}

	shutdown := otel.InitTracer("service-b")
	defer shutdown(context.Background())

	port := "8081"
	log.Printf("Service B running on port %s...", port)
	http.HandleFunc("/weather", handlers.HandleWeather)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
