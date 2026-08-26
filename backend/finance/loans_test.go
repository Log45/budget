package finance

import (
	"testing"
	"time"

	"Log45/budget/backend/models"
)

func TestGenerateScheduleAmortizesBalance(t *testing.T) {
	loan := models.Loan{
		Principal:      120000,
		CurrentBalance: 120000,
		Rate:           0.12,
		Term:           12,
		StartDate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	schedule := GenerateSchedule(loan, CalculateMonthlyPayment(loan))
	if len(schedule) != 12 {
		t.Fatalf("schedule length = %d, want 12", len(schedule))
	}
	if schedule[0].Interest <= schedule[len(schedule)-1].Interest {
		t.Errorf("interest did not decline: first=%d last=%d", schedule[0].Interest, schedule[len(schedule)-1].Interest)
	}
	if schedule[len(schedule)-1].Balance != 0 {
		t.Errorf("final balance = %d, want 0", schedule[len(schedule)-1].Balance)
	}
}
