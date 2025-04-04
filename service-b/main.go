package main

import (
	"context"
	"log"
	"net/http"

	"service-b/handlers"
	"service-b/otel"
)

func main() {
	shutdown := otel.InitTracer("service-b")
	defer shutdown(context.Background())

	port := "8081"
	log.Printf("Service B running on port %s...", port)
	http.HandleFunc("/weather", handlers.HandleWeather)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
