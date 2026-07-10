package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"Log45/budget/backend/api"
	"Log45/budget/backend/db"
	"Log45/budget/backend/services"
)

func main() {
	ctx := context.Background()

	pool, err := db.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepo := db.NewUserRepository(pool)

	userService := services.NewUserService(userRepo)

	authService := services.NewAuthService(
		userRepo,
		os.Getenv("JWT_SECRET"),
	)

	handler := api.NewHandler(*authService, *userService)

	router := api.RegisterRoutes()

	_ = handler // placeholder to remove errors, update functionality.

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
