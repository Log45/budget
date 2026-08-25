package api

import (
	"net/http"

	"Log45/budget/backend/services"
)

// Handler holds service dependencies used by HTTP route handlers.
type Handler struct {
	Auth  services.AuthService
	Users services.UserService
}

// NewHandler constructs a Handler with the services required by the API layer.
func NewHandler(auth services.AuthService, users services.UserService) Handler {
	return Handler{
		Auth:  auth,
		Users: users,
	}
}

// HealthHandler responds with 200 OK when the server process is running.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.Auth.Register(r.Context(), username, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := h.Auth.Login(r.Context(), username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

// LoanPaymentHandler is a placeholder for loan payment processing (Phase 3).
func (h *Handler) LoanPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := AuthenticatedUserID(r.Context()); !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Loan payment processed"))
}

// CreateBudgetHandler is a placeholder for budget creation (Phase 4).
func (h *Handler) CreateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := AuthenticatedUserID(r.Context()); !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Budget created"))
}

// GetBudgetHandler is a placeholder for budget retrieval (Phase 4).
func (h *Handler) GetBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := AuthenticatedUserID(r.Context()); !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Budget retrieved"))
}
