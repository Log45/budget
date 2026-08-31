package models

import "time"

type AccountType string

const (
	SavingsAccount      AccountType = "savings"
	CheckingAccount     AccountType = "checking"
	HybridAccount       AccountType = "hybrid"
	IndividualBrokerage AccountType = "individual_brokerage"
	RothIRAAccount      AccountType = "roth_ira"
	IRAAccount          AccountType = "ira"
	Account401K         AccountType = "401k"
	Account403B         AccountType = "403b"
	Account457B         AccountType = "457b"
	CreditCardAccount   AccountType = "credit_card"
)

type Account struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	Name        string      `json:"name"`
	AccountType AccountType `json:"account_type"`
	Balance     Money       `json:"balance"`
	Limit       *Money      `json:"limit,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
