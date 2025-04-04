package main

import (
	"context"
	"log"
	"net/http"
	"service-a/otel"
	"service-a/routes"
)

func main() {
	shutdown := otel.InitTracer("service-a")
	defer shutdown(context.Background())

	port := "8080"
	log.Printf("Service running on port %s...", port)
	http.ListenAndServe(":"+port, routes.NewRouter())
}
