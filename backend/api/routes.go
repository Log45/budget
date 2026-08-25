package api

import "net/http"

// RegisterRoutes wires all HTTP routes to the given handler and returns the mux.
func RegisterRoutes(h Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", HealthHandler)

	// Authentication routes
	mux.HandleFunc("POST /register", h.RegisterHandler)
	mux.HandleFunc("POST /login", h.LoginHandler)

	// User routes
	// TODO: Re-evaluate the need for user routes, if we have login, then we might not really need to query for the user except for maybe a user page?

	mux.HandleFunc("POST /loans/payment", h.LoanPaymentHandler)
	mux.HandleFunc("POST /budgets", h.CreateBudgetHandler)
	mux.HandleFunc("GET /budgets", h.GetBudgetHandler)

	return mux
}
