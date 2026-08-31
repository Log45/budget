package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"Log45/budget/backend/api"
	"Log45/budget/backend/db"
	"Log45/budget/backend/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	userRepo := db.NewUserRepository(pool)
	loanRepo := db.NewLoanRepository(pool)
	transactionRepo := db.NewTransactionRepository(pool)
	categoryRepo := db.NewCategoryRepository(pool)
	budgetRepo := db.NewBudgetRepository(pool)
	propertyRepo := db.NewPropertyRepository(pool)
	accountRepo := db.NewAccountRepository(pool)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, jwtSecret)
	loanService := services.NewLoanService(loanRepo)
	transactionService := services.NewTransactionService(transactionRepo, categoryRepo, budgetRepo, propertyRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	budgetService := services.NewBudgetService(budgetRepo, categoryRepo)
	propertyService := services.NewPropertyService(propertyRepo, loanRepo)
	accountService := services.NewAccountService(accountRepo)

	handler := api.NewHandler(*authService, *userService, loanService, transactionService, categoryService, budgetService, propertyService, accountService)
	router := api.RegisterRoutes(handler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
