package api

import "net/http"

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", HealthHandler)

	mux.HandleFunc("POST /loan/payment", LoanPaymentHandler)

	mux.HandleFunc("POST /budget", CreateBudgetHandler)

	mux.HandleFunc("GET /budget", GetBudgetHandler)

	return mux
}
