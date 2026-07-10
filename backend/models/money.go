package models

type Money int64

const (
	Dollar Money = 100 // 1 dollar = 100 cents, use int for balance to avoid floating point precision errors
)
