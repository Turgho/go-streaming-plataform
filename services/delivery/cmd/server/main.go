package main

import (
	"delivery-service/internal/handler"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		log.Fatal("failed to get server http port")
	}

	dirPath := os.Getenv("VIDEOS_PATH")
	if dirPath == "" {
		log.Fatal("failed to get videos dirPath")
	}

	mux := http.NewServeMux()
	h := handler.NewVideoHandler(dirPath)
	r := handler.NewVideoRoute(mux)
	r.RegisterRoutes(h)

	log.Printf("delivery-service rodando na porta :%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), handler.WithHeaders(mux))
}
