package main

import (
	"context"
	"go-microservice/internal/auth"
	"log"
	"net/http"
)

func main() {
	// in memory repository
	repo := auth.NewUserRepository()
	service := auth.NewAuthService(repo)
	ctx := context.Background()

	u1, err := service.Register(ctx, "test@example.com", "any_password")
	log.Println(u1, err)

	// hardcoded user registration for testing and demonstration purposes
	u2, err2 := service.Register(ctx, "test@example_2.com", "any_password")
	log.Println(u2, err2)

	if err2 != nil {
		log.Printf("Warning: test user registration failed: %v", err)
	}

	authHandler := auth.NewHandler(ctx, service)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", authHandler.Login)
	log.Println("Starting server on :8081...")

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
