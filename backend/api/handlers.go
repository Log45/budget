package api

import (
	"net/http"

	"Log45/budget/backend/services"
)

type Handler struct {
	Auth  services.AuthService
	Users services.UserService
	// TODO: Add services
}

func NewHandler(auth services.AuthService, users services.UserService) Handler

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func LoanPaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Handle loan payment logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Loan payment processed"))
}

func CreateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	// Handle budget creation logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Budget created"))
}

func GetBudgetHandler(w http.ResponseWriter, r *http.Request) {
	// Handle budget retrieval logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Budget retrieved"))
}
