package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

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
	Budgets      *services.BudgetService
	Properties   *services.PropertyService
	Accounts     *services.AccountService
}

// NewHandler constructs a Handler with the services required by the API layer.
func NewHandler(auth services.AuthService, users services.UserService, loans *services.LoanService, transactions *services.TransactionService, categories *services.CategoryService, budgets *services.BudgetService, properties *services.PropertyService, accounts *services.AccountService) Handler {
	return Handler{
		Auth:         auth,
		Users:        users,
		Loans:        loans,
		Transactions: transactions,
		Categories:   categories,
		Budgets:      budgets,
		Properties:   properties,
		Accounts:     accounts,
	}
}

func accountID(r *http.Request) (int64, error) { return strconv.ParseInt(r.PathValue("id"), 10, 64) }
func (h *Handler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	items, err := h.Accounts.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (h *Handler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Account
	if json.NewDecoder(r.Body).Decode(&item) != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	created, err := h.Accounts.Create(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := accountID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	item, err := h.Accounts.Get(r.Context(), userID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (h *Handler) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := accountID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	var item models.Account
	if json.NewDecoder(r.Body).Decode(&item) != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	item.ID = id
	userID, _ := AuthenticatedUserID(r.Context())
	updated, err := h.Accounts.Update(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
func (h *Handler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := accountID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	if err := h.Accounts.Delete(r.Context(), userID, id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
	if errors.Is(err, db.ErrBudgetNotFound) || errors.Is(err, db.ErrPropertyNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if errors.Is(err, db.ErrAccountNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if errors.Is(err, services.ErrInvalidBudget) || errors.Is(err, services.ErrInvalidProperty) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, services.ErrInvalidAccount) {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	item, err := decodeTransaction(r)
	if err != nil {
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
	item, err := decodeTransaction(r)
	if err != nil {
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

// transactionInput accepts the date-only values emitted by HTML date controls
// as well as RFC 3339 timestamps from API clients.
type transactionInput struct {
	BudgetID      *int64                 `json:"budget_id"`
	CategoryID    *int64                 `json:"category_id"`
	PropertyID    *int64                 `json:"property_id"`
	AccountID     *int64                 `json:"account_id"`
	Amount        models.Money           `json:"amount"`
	Description   string                 `json:"description"`
	Type          models.TransactionType `json:"type"`
	Source        string                 `json:"source"`
	Destination   string                 `json:"destination"`
	Date          string                 `json:"date"`
	Recurring     bool                   `json:"recurring"`
	RecurringType models.BudgetType      `json:"recurring_type"`
}

func decodeTransaction(r *http.Request) (models.Transaction, error) {
	var input transactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return models.Transaction{}, err
	}
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		date, err = time.Parse(time.RFC3339, input.Date)
		if err != nil {
			return models.Transaction{}, err
		}
	}
	return models.Transaction{BudgetID: input.BudgetID, CategoryID: input.CategoryID, PropertyID: input.PropertyID, AccountID: input.AccountID, Amount: input.Amount, Description: input.Description, Type: input.Type, Source: input.Source, Destination: input.Destination, Date: date, Recurring: input.Recurring, RecurringType: input.RecurringType}, nil
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
	loan, err := decodeLoan(r)
	if err != nil {
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

type loanInput struct {
	Name           string       `json:"name"`
	Principal      models.Money `json:"principal"`
	CurrentBalance models.Money `json:"current_balance"`
	Rate           float64      `json:"rate"`
	Term           int          `json:"term"`
	StartDate      string       `json:"start_date"`
}

func decodeLoan(r *http.Request) (models.Loan, error) {
	var input loanInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return models.Loan{}, err
	}
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		startDate, err = time.Parse(time.RFC3339, input.StartDate)
		if err != nil {
			return models.Loan{}, err
		}
	}
	return models.Loan{Name: input.Name, Principal: input.Principal, CurrentBalance: input.CurrentBalance, Rate: input.Rate, Term: input.Term, StartDate: startDate}, nil
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

func (h *Handler) CreateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Budget
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	created, err := h.Budgets.Create(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) ListBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	items, err := h.Budgets.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func budgetID(r *http.Request) (int64, error) { return strconv.ParseInt(r.PathValue("id"), 10, 64) }
func (h *Handler) GetBudgetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := budgetID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid budget id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	item, err := h.Budgets.Get(r.Context(), userID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (h *Handler) UpdateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := budgetID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid budget id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Budget
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	item.ID = id
	updated, err := h.Budgets.Update(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
func (h *Handler) DeleteBudgetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := budgetID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid budget id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	if err := h.Budgets.Delete(r.Context(), userID, id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func propertyID(r *http.Request) (int64, error) { return strconv.ParseInt(r.PathValue("id"), 10, 64) }
func (h *Handler) CreatePropertyHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Property
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	created, err := h.Properties.Create(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) ListPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := AuthenticatedUserID(r.Context())
	items, err := h.Properties.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (h *Handler) GetPropertyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := propertyID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	item, err := h.Properties.Get(r.Context(), userID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (h *Handler) UpdatePropertyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := propertyID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	var item models.Property
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	item.ID = id
	updated, err := h.Properties.Update(r.Context(), userID, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
func (h *Handler) DeletePropertyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := propertyID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	if err := h.Properties.Delete(r.Context(), userID, id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) PropertyAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := propertyID(r)
	if err != nil || id <= 0 {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}
	userID, _ := AuthenticatedUserID(r.Context())
	item, err := h.Properties.Analytics(r.Context(), userID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
