# Dev Notes

API Examples:

Method	Endpoint	Service
POST	/users	UserService
GET	/users/{id}	UserService
GET	/transactions	TransactionService
POST	/transactions	TransactionService
GET	/budgets	BudgetService
POST	/budgets	BudgetService
GET	/loans	LoanService
POST	/loans/payment	LoanService
POST	/taxes/calculate	TaxService
GET	/investments/forecast	InvestmentService
POST	/plaid/link	PlaidService

Each endpoint maps cleanly to a handler, which delegates to a corresponding service.






The responsibility of each layer is:

api: Parse requests and return responses.
services: Registration, login, password hashing, JWT generation, authorization.
db: Read and write records.
models: Domain objects.