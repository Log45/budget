package finance

import (
	"Log45/budget/backend/models"
	"math"
)

func MonthlyInterest(balance models.Money, rate float64) models.Money {
	return models.Money(float64(balance) * rate / 12)
}

func AnnualInterest(balance models.Money, rate float64) models.Money {
	return models.Money(float64(balance) * rate)
}

func CompoundInterest(balance models.Money, rate float64, periods int, years int) models.Money {
	return models.Money(float64(balance) * math.Pow(1+rate/float64(periods), float64(periods*years)))
}
