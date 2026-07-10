package finance

import (
	"Log45/budget/backend/models"
	"errors"
)

var TransactionTypeError = errors.New("Invalid transaction type")

func CalculateBudgetRatio(transaction models.Transaction, netPay models.Money) (float64, error) {
	if transaction.Type != models.Expense {
		return 0, TransactionTypeError
	}
	if netPay == 0 {
		return 0, errors.New("netPay must be greater than 0")
	}
	expense := transaction.Amount
	ratio := float64(expense) / float64(netPay)
	return ratio, nil
}
