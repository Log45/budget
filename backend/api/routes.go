package api

import "net/http"

// RegisterRoutes wires all HTTP routes to the given handler and returns the mux.
func RegisterRoutes(h Handler) http.Handler {
	mux := http.NewServeMux()
	requireAuth := AuthMiddleware(h.Auth)

	mux.HandleFunc("GET /health", HealthHandler)

	// Authentication routes
	mux.HandleFunc("POST /register", h.RegisterHandler)
	mux.HandleFunc("POST /login", h.LoginHandler)
	// /auth paths are the public API names; retain the original paths for the
	// existing frontend proxy.
	mux.HandleFunc("POST /auth/register", h.RegisterHandler)
	mux.HandleFunc("POST /auth/login", h.LoginHandler)

	// A profile is always scoped to its authenticated owner.
	mux.Handle("GET /users/{id}", requireAuth(http.HandlerFunc(h.GetUserHandler)))

	mux.Handle("GET /loans", requireAuth(http.HandlerFunc(h.ListLoansHandler)))
	mux.Handle("POST /loans", requireAuth(http.HandlerFunc(h.CreateLoanHandler)))
	mux.Handle("GET /loans/{id}", requireAuth(http.HandlerFunc(h.GetLoanHandler)))
	mux.Handle("GET /loans/{id}/schedule", requireAuth(http.HandlerFunc(h.LoanScheduleHandler)))
	mux.Handle("GET /loans/{id}/analytics", requireAuth(http.HandlerFunc(h.LoanAnalyticsHandler)))
	mux.Handle("POST /budgets", requireAuth(http.HandlerFunc(h.CreateBudgetHandler)))
	mux.Handle("GET /budgets", requireAuth(http.HandlerFunc(h.ListBudgetsHandler)))
	mux.Handle("GET /budgets/{id}", requireAuth(http.HandlerFunc(h.GetBudgetHandler)))
	mux.Handle("PUT /budgets/{id}", requireAuth(http.HandlerFunc(h.UpdateBudgetHandler)))
	mux.Handle("DELETE /budgets/{id}", requireAuth(http.HandlerFunc(h.DeleteBudgetHandler)))
	mux.Handle("GET /categories", requireAuth(http.HandlerFunc(h.ListCategoriesHandler)))
	mux.Handle("POST /categories", requireAuth(http.HandlerFunc(h.CreateCategoryHandler)))
	mux.Handle("GET /transactions", requireAuth(http.HandlerFunc(h.ListTransactionsHandler)))
	mux.Handle("POST /transactions", requireAuth(http.HandlerFunc(h.CreateTransactionHandler)))
	mux.Handle("PUT /transactions/{id}", requireAuth(http.HandlerFunc(h.UpdateTransactionHandler)))
	mux.Handle("DELETE /transactions/{id}", requireAuth(http.HandlerFunc(h.DeleteTransactionHandler)))
	mux.Handle("GET /properties", requireAuth(http.HandlerFunc(h.ListPropertiesHandler)))
	mux.Handle("POST /properties", requireAuth(http.HandlerFunc(h.CreatePropertyHandler)))
	mux.Handle("GET /properties/{id}", requireAuth(http.HandlerFunc(h.GetPropertyHandler)))
	mux.Handle("PUT /properties/{id}", requireAuth(http.HandlerFunc(h.UpdatePropertyHandler)))
	mux.Handle("DELETE /properties/{id}", requireAuth(http.HandlerFunc(h.DeletePropertyHandler)))
	mux.Handle("GET /properties/{id}/analytics", requireAuth(http.HandlerFunc(h.PropertyAnalyticsHandler)))

	return mux
}
