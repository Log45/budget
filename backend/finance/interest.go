package finance

import "math"

func MonthlyInterest(balance Money, rate float64) Money {
	return Money(float64(balance) * rate / 12)
}

func AnnualInterest(balance Money, rate float64) Money {
	return Money(float64(balance) * rate)
}

func CompoundInterest(balance Money, rate float64, periods int, years int) Money {
	return Money(float64(balance) * math.Pow(1+rate/float64(periods), float64(periods*years)))
}
