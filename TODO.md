# Budget App — Prototype Todo List

This tracks work needed before the budgeting/finances website prototype is runnable with the target feature set.

**Current state:** Backend-only Go scaffold (~31 files). Some finance math and auth logic exist, but the app did not compile, had no database schema, no frontend, and no Properties feature.

---

## Phase 0 — Unblock running locally

- [x] **Implement `UserRepository` CRUD** — all 6 methods in `backend/db/user_repository.go`
- [x] **Complete stub methods** — `NewHandler` in `backend/api/handlers.go`; `GetByUsername`, `Update`, `Delete` in `backend/services/users.go`
- [x] **Add database schema / migrations** — initial SQL migration for users, loans, budgets, transactions, accounts, properties
- [x] **Add `.env.example` and setup docs** — document `DATABASE_URL`, `JWT_SECRET`, Plaid vars; add root README
- [x] **Wire dependency injection in `main.go`** — connect services → handlers → routes
- [x] **Verify `go build` and `go run .` succeed** with Postgres via `docker compose up -d`
- [x] **Align API paths with `docs/dev.md`** — `/budgets`, `/loans/payment`

---

## Phase 1 — Auth & user foundation

- [x] **Auth HTTP routes** — `POST /auth/register`, `POST /auth/login`
- [ ] **JWT auth middleware** — protect all user-scoped endpoints 
- [ ] **User routes** — `POST /users`, `GET /users/{id}` per `docs/dev.md`
- [ ] **Subscription / tier model** — distinguish free (manual entry, simulation) vs paid (Plaid sync)

---

## Phase 2 — Frontend (does not exist yet)

- [x] **Choose and scaffold frontend** — React/Next.js or similar with tab-based navigation
- [x] **Auth UI** — register, login, session/token handling
- [x] **API client layer** — typed calls to backend endpoints
- [x] **Shared layout** — nav tabs: Budget, Expenses, Properties, Loans, (Investments)
- [x] **Dev proxy / CORS** — frontend on one port talking to backend `:8080`

---

## Phase 3 — Loans tab

**Requirement:** Visualize loans, balance & interest, predicted future balance, payment scheduling, interest vs. principal analytics.

- [ ] **Revisit Loan logic/required features** - Before implementing API endpoints, it needs to be determined what is needed for proper analysis of loans. 
- [ ] **Fix amortization schedule** — `GenerateSchedules` needs per-payment principal/interest/balance breakdown
- [ ] **Add loan payment simulation** — apply extra payments, recalculate schedule, show total interest saved
- [ ] **Loan persistence** — implement `LoanRepository` (currently struct-only stub)
- [ ] **Loan service** — CRUD, payment recording, schedule generation, future balance projection
- [ ] **Loan API routes** — `GET /loans`, `POST /loans`, `GET /loans/{id}`, `POST /loans/payment`, `GET /loans/{id}/schedule`, `GET /loans/{id}/analytics`
- [ ] **Loan analytics endpoints** — cumulative interest vs. principal, payoff date, remaining interest
- [ ] **Loans tab UI** — loan list, detail view, amortization chart/table, payment simulator, interest vs. principal visualization

---

## Phase 4 — Monthly budget (salary + tax/withholdings)

**Requirement:** Monthly budget based on salary with state + federal income tax and withholdings.

- [ ] **Tax/withholding engine** — implement or integrate federal + state brackets, FICA, standard deduction, filing status
- [ ] **Withholding inputs model** — salary, pay frequency, filing status, state, W-4 allowances/extra withholding, pre-tax deductions
- [ ] **Net pay calculation** — gross → federal tax → state tax → FICA → net monthly/periodic income
- [ ] **Budget service** — create monthly budget from net pay; allocate categories; track spent vs. planned
- [ ] **Budget persistence** — implement `BudgetRepository`
- [ ] **Budget API routes** — `GET /budgets`, `POST /budgets`, `GET /budgets/{id}`, update categories
- [ ] **Budget tab UI** — salary/tax inputs, net pay display, category breakdown, spending vs. budget bars

---

## Phase 5 — Expenses tab

**Requirement:** Imported from banks or manually entered.

- [ ] **Account model** — replace TODO on `Transaction.Source`/`Destination` with proper `Account` entity
- [ ] **Transaction repository & service** — CRUD, filtering by date/category/account
- [ ] **Expense categorization** — categories/tags, optional auto-categorization rules
- [ ] **Transaction API routes** — `GET /transactions`, `POST /transactions`, `PUT /transactions/{id}`, `DELETE /transactions/{id}`
- [ ] **Manual entry UI** — form for amount, date, category, account, description
- [ ] **Expenses tab UI** — list/filter by month, category totals, import vs. manual indicator
- [ ] **Link expenses to budget** — roll up spending against budget categories

---

## Phase 6 — Plaid integration (paid tier)

**Requirement:** Paid product connects to banks, cards, and investments via Plaid.

- [ ] **Plaid env config** — use `PLAID_ENV` (sandbox for dev); fix hardcoded `plaid.Production`
- [ ] **Plaid Link flow** — `POST /plaid/link/token`, `POST /plaid/link/exchange`
- [ ] **Store Plaid items & accounts** — DB tables for linked institutions, account mappings
- [ ] **Transaction sync** — pull transactions from Plaid, dedupe, map to internal `Transaction` model
- [ ] **Investment sync** — pull holdings/balances for investment accounts
- [ ] **Webhook handler** — Plaid transaction updates
- [ ] **Tier gating** — only paid users can initiate Plaid Link
- [ ] **Plaid Link UI** — connect bank/card/investment buttons, account picker, sync status
- [ ] **Plaid sandbox testing** — test credentials, sandbox institution flow

---

## Phase 7 — Properties tab (entirely missing)

**Requirement:** Reference mortgage loan, value, utilities, rent/income, renovation/repair expenses; optionally apartments for rent tracking.

- [ ] **Property model** — address, type (house, apartment, rental), purchase price, current value, linked mortgage loan ID
- [ ] **Apartment / rental unit model** — track rent expenses vs. rental income
- [ ] **Property expense types** — utilities, insurance, HOA, maintenance, repairs, renovations, property tax
- [ ] **Property income** — rent received, vacancy, other income
- [ ] **Link property ↔ loan** — associate a property with its mortgage from the Loans tab
- [ ] **Property repository & service**
- [ ] **Property API routes** — CRUD properties, expenses, income, utilities
- [ ] **Property analytics** — cash flow, equity, cap rate for rentals
- [ ] **Properties tab UI** — property cards, linked loan summary, expense/income ledger, equity chart
- [ ] **Apartment rent expense tracking** — separate view or sub-type for personal rent payments

---

## Phase 8 — Cross-cutting & prototype polish

- [ ] **Consistent API response format** — JSON envelopes, error codes
- [ ] **Input validation** — request parsing on all handlers
- [ ] **Money handling audit** — ensure JSON serialization is consistent (dollars vs. cents)
- [ ] **Seed data script** — sample user, salary, loan, property, transactions for demo
- [ ] **Dockerfile for backend** — uncomment and finish backend service in `docker-compose.yml`
- [ ] **Basic tests** — at minimum for loan amortization, tax calc, and budget ratio
- [ ] **Health check expansion** — DB connectivity check on `GET /health`

---

## Suggested build order

| Order | Work | Why |
|-------|------|-----|
| 1 | Phase 0 (compile + DB + wiring) | Nothing runs without it |
| 2 | Phase 1 (auth) | Multi-user data isolation |
| 3 | Phase 2 (frontend shell + tabs) | UI needed for prototype |
| 4 | Phase 3 (loans) | Most backend math already started |
| 5 | Phase 4 (budget + taxes) | Core value prop |
| 6 | Phase 5 (manual expenses) | Works without Plaid |
| 7 | Phase 7 (properties) | New domain, self-contained |
| 8 | Phase 6 (Plaid) | Paid tier; defer for free-tier prototype |
| 9 | Phase 8 (polish) | Demo readiness |

---

## Feature status matrix

| Feature | Backend model | Finance logic | DB/repo | API | Frontend |
|---------|--------------|---------------|---------|-----|----------|
| Loans | ✅ | ⚠️ partial | ❌ | ⚠️ stub | ❌ |
| Budget | ✅ | ⚠️ ratio only | ❌ | ⚠️ stub | ❌ |
| Expenses | ✅ | ⚠️ ratio only | ❌ | ❌ | ❌ |
| Taxes | ❌ | ❌ | ❌ | ❌ | ❌ |
| Plaid | ❌ | ❌ | ❌ | ❌ | ❌ |
| Properties | ❌ | ❌ | ⚠️ schema only | ❌ | ❌ |
| Auth | ✅ logic | — | ✅ | ❌ | ❌ |
| Investments | ✅ model | ❌ | ❌ | ❌ | ❌ |
