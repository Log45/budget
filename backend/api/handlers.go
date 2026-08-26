package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
	"Log45/budget/backend/services"
)

// Handler holds service dependencies used by HTTP route handlers.
type Handler struct {
	Auth         services.AuthService
	Users        services.UserService
	Loans        *services.LoanService
	Transactions *services.TransactionService
	Categories   *services.CategoryService
}

// NewHandler constructs a Handler with the services required by the API layer.
func NewHandler(auth services.AuthService, users services.UserService, loans *services.LoanService, transactions *services.TransactionService, categories *services.CategoryService) Handler {
	return Handler{
		Auth:         auth,
		Users:        users,
		Loans:        loans,
		Transactions: transactions,
		Categories:   categories,
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrLoanNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if errors.Is(err, services.ErrInvalidLoan) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, services.ErrInvalidTransaction) || errors.Is(err, services.ErrInvalidCategory) || errors.Is(err, db.ErrCategoryNotFound) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, db.ErrTransactionNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func (h *Handler) ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	items, err := h.Categories.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (h *Handler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Category
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	created, err := h.Categories.Create(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	items, err := h.Transactions.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (h *Handler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	created, err := h.Transactions.Create(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	item.ID = r.PathValue("id")
	updated, err := h.Transactions.Update(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
func (h *Handler) DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	if err := h.Transactions.Delete(r.Context(), userID, r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || requestedID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	if requestedID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	user, err := h.Users.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) CreateLoanHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	var loan models.Loan
	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	loan.UserID = userID
	created, err := h.Loans.Create(r.Context(), loan)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) ListLoansHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	loans, err := h.Loans.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, loans)
}

func loanID(r *http.Request) (int64, error) { return strconv.ParseInt(r.PathValue("id"), 10, 64) }
func (h *Handler) GetLoanHandler(w http.ResponseWriter, r *http.Request) {
	id, err := loanID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid loan id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	loan, err := h.Loans.Get(r.Context(), userID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, loan)
}
func (h *Handler) LoanScheduleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := loanID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid loan id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	schedule, err := h.Loans.Schedule(r.Context(), userID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, schedule)
}
func (h *Handler) LoanAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := loanID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid loan id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	analytics, err := h.Loans.Analytics(r.Context(), userID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, analytics)
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
